package gdrive

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vasyahuyasa/july/opds/storage"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	MIMETypeFolder = "application/vnd.google-apps.folder"
)

var _ storage.Storage = &GdriveStorage{}

type GdriveStorage struct {
	rootID string
	svc    *drive.Service
}

func NewStorage(root string, apikey string) (*GdriveStorage, error) {
	ctx := context.Background()

	svc, err := drive.NewService(ctx, option.WithAPIKey(apikey), option.WithScopes(drive.DriveMetadataReadonlyScope))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %w", err)
	}

	return &GdriveStorage{
		rootID: root,
		svc:    svc,
	}, nil
}

func (drive *GdriveStorage) List(path string) ([]storage.StorageEntry, error) {
	fileID := drive.filterpath(path)
	log.Println("List:", fileID)

	var entries []storage.StorageEntry

	nextPageToken := ""
	for {
		r, err := drive.svc.Files.List().
			Q(fmt.Sprintf("%q in parents", fileID)).
			PageToken(nextPageToken).
			Fields("nextPageToken, files(*)").
			Do()
		if err != nil {
			return nil, fmt.Errorf("unable to list files in %q: %w", fileID, err)
		}

		for _, i := range r.Files {
			fmt.Println(i.Name, i.Id, i.MimeType)

			mod, err := time.Parse(time.RFC822, i.ModifiedTime)
			if err != nil {
				return nil, err
			}

			entries = append(entries, storage.StorageEntry{
				Title:    i.Name,
				Path:     i.Id,
				IsDir:    i.MimeType == MIMETypeFolder,
				Updated:  mod,
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
	fileID := drive.filterpath(path)

	f, err := drive.svc.Files.Get(fileID).Do()
	if err != nil {
		return false, fmt.Errorf("unable to get file %q: %w", fileID, err)
	}

	fmt.Println("GET", f.Name)
	return f.MimeType != MIMETypeFolder, nil
}

func (drive *GdriveStorage) Download(w io.Writer, path string) error {
	_ = drive.filterpath(path)
	panic("not implemented")
}

func (drive *GdriveStorage) filterpath(path string) string {
	if path == "" {
		return drive.rootID
	}

	return path
}
