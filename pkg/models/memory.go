package models

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpeg "github.com/dsoprea/go-jpeg-image-structure/v2"
)

type Memory struct {
	Date             string `json:"Date"`
	MediaType        string `json:"Media Type"`
	Location         string `json:"Location"`
	DownloadLink     string `json:"Download Link"`
	MediaDownloadUrl string `json:"Media Download Url"`

	userId   string
	uniqueId string
	fileName string
}

type Memories struct {
	SavedMedia []Memory `json:"Saved Media"`
}

func (m *Memory) Save(saveDir string) error {
	saveFile, err := m.FileName()
	if err != nil {
		return err
	}

	outputFile := filepath.Join(saveDir, saveFile)
	// Download file

	sleepTime := 5 * time.Second
	success := false
	for i := 1; i <= 10; i++ {

		err = _saveRemoteFile(m.MediaDownloadUrl, outputFile)

		if err != nil {
			// we should retry
			//log.Println("attempt failed", i, err)
			time.Sleep(sleepTime)
			sleepTime = time.Duration(5*i) * time.Second
			continue
		} else {
			//if i > 1 {
			//log.Println("attempt succeeded", i)
			//}
			success = true
			break
		}
	}

	if !success {
		return errors.New(fmt.Sprintf("Failed to save to disk: %s", err))
	}

	newFileName, err := _renameToDetectedType(outputFile)

	if err != nil {
		return err
	}

	ts, err := time.Parse("2006-01-02 15:04:05 UTC", m.Date)
	if err != nil {
		return err
	}

	//if m.MediaType == "Image" && strings.HasSuffix(newFileName, ".zip") {
	// handle jpgs with pngs
	//}

	if m.MediaType == "Image" && strings.HasSuffix(newFileName, ".jpg") {
		// images dont have exif date for when it was created
		return SetDateIfNone(newFileName, ts)
	}

	return nil
}

func (m *Memory) UserId() (string, error) {
	if m.userId == "" {

		u, err := url.Parse(m.DownloadLink)
		if err != nil {
			return "", err
		}

		q, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			return "", err
		}

		return q.Get("uid"), nil
	}

	return m.userId, nil
}

func (m *Memory) UniqueId() (string, error) {
	if m.uniqueId == "" {

		u, err := url.Parse(m.MediaDownloadUrl)
		if err != nil {
			return "", err
		}

		q, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			return "", err
		}

		return q.Get("sig"), nil
	}
	return m.uniqueId, nil
}

func (m *Memory) FileName() (string, error) {
	if m.fileName == "" {

		mid, err := m.UniqueId()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s_%s", m.Date, mid), nil
	}

	return m.fileName, nil
}

func _saveRemoteFile(urlToSave, outputFileName string) error {
	// Create the HTTP client
	client := &http.Client{}

	// Create the GET request
	req, err := http.NewRequest("GET", urlToSave, nil)
	if err != nil {
		return err
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("invalid response code %d, failed to get %s", resp.StatusCode, urlToSave))
	}

	// Create the output file
	outFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outFile.Close() // Ensure the output file is closed

	// Stream the response body to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	_, err = outFile.Seek(0, 0)
	if err != nil {
		return err
	}

	return nil
}

func _renameToDetectedType(outputFileName string) (string, error) {
	fileDetect, err := os.Open(outputFileName)
	defer fileDetect.Close()

	if err != nil {
		return outputFileName, err
	}
	// Read the first 512 bytes for content sniffing
	buffer := make([]byte, 512)
	_, err = fileDetect.Read(buffer)
	if err != nil && err != io.EOF {
		return outputFileName, err
	}

	contentType := http.DetectContentType(buffer)
	newFileName := outputFileName
	if contentType == "image/jpeg" {
		newFileName = outputFileName + ".jpg"
	} else if contentType == "video/mp4" {
		newFileName = outputFileName + ".mp4"
	} else if contentType == "application/zip" {
		newFileName = outputFileName + ".zip"
	} else {
		newFileName = outputFileName + ".mov"
	}
	err = os.Rename(outputFileName, newFileName)
	return newFileName, err
}

const (
	tagDateTimeOriginal = "DateTimeOriginal"
	tagDateTime         = "DateTime"
)

func setTag(rootIB *exif.IfdBuilder, ifdPath, tagName, tagValue string) error {
	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIB, ifdPath)
	if err != nil {
		return fmt.Errorf("Failed to get or create IB: %v", err)
	}

	if err := ifdIb.SetStandardWithName(tagName, tagValue); err != nil {
		return fmt.Errorf("failed to set DateTime tag: %v", err)
	}

	return nil
}

// SetDateIfNone sets the EXIF DateTime to the given Time unless it has
// already been defined.
func SetDateIfNone(filepath string, t time.Time) error {
	parser := jpeg.NewJpegMediaParser()
	psl, err := parser.ParseFile(filepath)
	if err != nil {
		return fmt.Errorf("Failed to parse JPEG file: %s - %v", filepath, err)
	}

	sl := psl.(*jpeg.SegmentList)

	rootIb, err := sl.ConstructExifBuilder()
	if err != nil {
		im, err := exifcommon.NewIfdMappingWithStandard()
		if err != nil {
			log.Fatal(err)
		}
		ti := exif.NewTagIndex()
		if err := exif.LoadStandardTags(ti); err != nil {
			return fmt.Errorf("Failed to load standard tags: %v", err)
		}

		rootIb = exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, exifcommon.EncodeDefaultByteOrder)
	}

	// Form our timestamp string
	ts := exifcommon.ExifFullTimestampString(t)

	// Set DateTime
	ifdPath := "IFD0"
	if err := setTag(rootIb, ifdPath, tagDateTime, ts); err != nil {
		return fmt.Errorf("Failed to set tag %v: %v", tagDateTime, err)
	}

	// Set DateTimeOriginal
	ifdPath = "IFD/Exif"
	if err := setTag(rootIb, ifdPath, tagDateTimeOriginal, ts); err != nil {
		return fmt.Errorf("Failed to set tag %v: %v", tagDateTime, err)
	}

	// Update the exif segment.
	if err := sl.SetExif(rootIb); err != nil {
		return fmt.Errorf("Failed to set EXIF to jpeg: %v", err)
	}

	// Write the modified file
	b := new(bytes.Buffer)
	if err := sl.Write(b); err != nil {
		return fmt.Errorf("Failed to create JPEG data: %v", err)
	}

	// Save the file
	if err := os.WriteFile(filepath, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("Failed to write JPEG file: %v", err)
	}

	return nil
}
