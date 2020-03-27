package gdrive

import (
	"io"

	"github.com/vasyahuyasa/july/opds/storage"
	"google.golang.org/api/drive/v3"
)

const (
	MIMETypeFolder = "application/vnd.google-apps.folder"
)

var _ storage.Storage = &GdriveStorage{}

type GdriveStorage struct {
	svc *drive.Service
}

func (drive *GdriveStorage) List(path string) ([]storage.StorageEntry, error) {
	var entries []storage.StorageEntry

	nextPageToken := ""
	for {
		r, err := drive.svc.Files.List().
			Q("'" + path + "' in parents").
			PageToken(nextPageToken).
			Fields("nextPageToken, files(*)").
			Do()
		if err != nil {
			return nil, err
		}

		for _, i := range r.Files {
			fmt.Println(i.Name, i.Id, i.MimeType)

			entries = append(entries, storage.StorageEntry{
				Title:    i.Name,
				Path:     i.Id,
				IsDir:    i.MimeType == MIMETypeFolder,
				Updated:  i.ModifiedTime,
				MimeType: i.MimeType,
			})
		}

		nextPageToken = r.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return entries, nil
}

func (drive *GdriveStorage) IsDownloadable(path string) (bool, error) {

}

func (drive *GdriveStorage) Download(w io.Writer, path string) error {

}
