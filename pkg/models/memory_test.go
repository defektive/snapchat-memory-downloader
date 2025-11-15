package models

import "testing"

func TestMemory_UserId(t *testing.T) {
	type fields struct {
		Date             string
		MediaType        string
		Location         string
		DownloadLink     string
		MediaDownloadUrl string
		userId           string
		uniqueId         string
		fileName         string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "userId",
			fields: fields{
				Date:             "2025-04-20 16:20:01 UTC",
				MediaType:        "Image",
				Location:         "-1,-1",
				DownloadLink:     "https://app.snapchat.com.example/dmd/memories?uid=uid-param-value&sid=sid-param-value&mid=mid-param-value&ts=ts-param-value&proxy=true&sig=download-sig-param-value",
				MediaDownloadUrl: "https://us-east1-aws.api.snapchat.com.example/dmd/mm?uid=uid-param-value&sid=sid-param-value&mid=mid-param-value&ts=ts-param-value&sig=media-sig-param-value",
			},
			want:    "uid-param-value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				Date:             tt.fields.Date,
				MediaType:        tt.fields.MediaType,
				Location:         tt.fields.Location,
				DownloadLink:     tt.fields.DownloadLink,
				MediaDownloadUrl: tt.fields.MediaDownloadUrl,
				userId:           tt.fields.userId,
				uniqueId:         tt.fields.uniqueId,
				fileName:         tt.fields.fileName,
			}
			got, err := m.UserId()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UserId() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemory_UniqueId(t *testing.T) {
	type fields struct {
		Date             string
		MediaType        string
		Location         string
		DownloadLink     string
		MediaDownloadUrl string
		userId           string
		uniqueId         string
		fileName         string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "userId",
			fields: fields{
				Date:             "2025-04-20 16:20:01 UTC",
				MediaType:        "Image",
				Location:         "-1,-1",
				DownloadLink:     "https://app.snapchat.com.example/dmd/memories?uid=uid-param-value&sid=sid-param-value&mid=mid-param-value&ts=ts-param-value&proxy=true&sig=download-sig-param-value",
				MediaDownloadUrl: "https://us-east1-aws.api.snapchat.com.example/dmd/mm?uid=uid-param-value&sid=sid-param-value&mid=mid-param-value&ts=ts-param-value&sig=media-sig-param-value",
			},
			want:    "media-sig-param-value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Memory{
				Date:             tt.fields.Date,
				MediaType:        tt.fields.MediaType,
				Location:         tt.fields.Location,
				DownloadLink:     tt.fields.DownloadLink,
				MediaDownloadUrl: tt.fields.MediaDownloadUrl,
				userId:           tt.fields.userId,
				uniqueId:         tt.fields.uniqueId,
				fileName:         tt.fields.fileName,
			}
			got, err := m.UniqueId()
			if (err != nil) != tt.wantErr {
				t.Errorf("UniqueId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UniqueId() got = %v, want %v", got, tt.want)
			}
		})
	}
}
