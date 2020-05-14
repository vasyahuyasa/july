package storage

import (
	"io"
	"time"
)

type StorageEntry struct {
	Title    string
	Path     string
	IsDir    bool
	Updated  time.Time
	MimeType string
}

type Storage interface {
	List(path string) ([]StorageEntry, error)
	IsDownloadable(path string) (bool, error)
	Download(w io.Writer, path string) error
}
