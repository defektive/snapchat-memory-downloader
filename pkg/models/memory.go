package models

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
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
			log.Println("attempt failed", i, err)
			time.Sleep(sleepTime)
			sleepTime = time.Duration(5*i) * time.Second
			continue
		} else {
			if i > 1 {
				log.Println("attempt succeeded", i)
			}
			success = true
			break
		}
	}

	if !success {
		return errors.New(fmt.Sprintf("Failed to save to disk: %s", err))
	}

	return _renameToDetectedType(outputFile)
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

func _renameToDetectedType(outputFileName string) error {
	fileDetect, err := os.Open(outputFileName)
	defer fileDetect.Close()

	if err != nil {
		return err
	}
	// Read the first 512 bytes for content sniffing
	buffer := make([]byte, 512)
	_, err = fileDetect.Read(buffer)
	if err != nil && err != io.EOF {
		return err
	}

	contentType := http.DetectContentType(buffer)
	if contentType == "image/jpeg" {
		err = os.Rename(outputFileName, outputFileName+".jpg")
	} else if contentType == "video/mp4" {
		err = os.Rename(outputFileName, outputFileName+".mp4")
	} else if contentType == "application/zip" {
		err = os.Rename(outputFileName, outputFileName+".zip")
	} else {
		err = os.Rename(outputFileName, outputFileName+".mov")
	}
	if err != nil {
		return err
	}

	return nil
}

func setExifDateTime(imagePath string, newTime time.Time) error {
	//
	//// Parse the JPEG structure
	//media, err := jpegstructure.NewJpegMediaParser().ParseFile(imagePath)
	//if err != nil {
	//	return fmt.Errorf("failed to parse JPEG structure: %w", err)
	//}
	//
	//
	//sl := media.(*jpegstructure.SegmentList)
	//// Get the root IFD
	//exifBuilder, err := sl.ConstructExifBuilder()
	//if err != nil {
	//	return fmt.Errorf("failed to get exifBuilder: %w", err)
	//}
	//
	//ifd0Ib, err := exif.GetOrCreateIbFromRootIb(exifBuilder, "IFD0")
	//if err != nil {
	//	return fmt.Errorf("failed to get IFD0 from root IB: %w", err)
	//}

	//ifd0Ib.SetStandardWithName("DateTimeOriginal", time.Now().)
	// Format the new time for EXIF
	//ts := newTime.Format("2006:01:02 15:04:05")

	//// Set DateTime (ModifyDate) in IFD0
	//if err := rootIb(exif.IfdPathStandard, exif.TagDateTime, ts); err != nil {
	//	return fmt.Errorf("failed to set DateTime tag: %w", err)
	//}
	//
	//// Set DateTimeOriginal in ExifIFD
	//exifIfd, err := rootIb.ChildExifIfd()
	//if err != nil {
	//	return fmt.Errorf("failed to get ExifIFD: %w", err)
	//}
	//if err := exifIfd.SetTag(exif.IfdPathStandard, exif.TagDateTimeOriginal, ts); err != nil {
	//	return fmt.Errorf("failed to set DateTimeOriginal tag: %w", err)
	//}
	//
	//// Set DateTimeDigitized in ExifIFD (optional)
	//if err := exifIfd.SetTag(exif.IfdPathStandard, exif.TagDateTimeDigitized, ts); err != nil {
	//	return fmt.Errorf("failed to set DateTimeDigitized tag: %w", err)
	//}
	//
	//// Write the modified image to a new file or overwrite the original
	//outputFile, err := os.Create("output_" + imagePath)
	//if err != nil {
	//	return fmt.Errorf("failed to create output file: %w", err)
	//}
	//defer outputFile.Close()
	//
	//if err := media.Write(outputFile); err != nil {
	//	return fmt.Errorf("failed to write modified JPEG: %w", err)
	//}

	return nil
}
