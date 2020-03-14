package local

import (
	"io"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/vasyahuyasa/july/opds/storage"
)

var _ storage.Storage = &FsStorage{}

var mimes = map[string]string{
	".epub": "application/epub+zip",
	".fb2":  "application/fb2+zip",
	".mobi": "application/x-mobipocket-ebook",
}

type FsStorage struct {
	root string
}

func NewFsStorage(root string) *FsStorage {
	return &FsStorage{
		root: root,
	}
}

func (store *FsStorage) List(path string) ([]storage.StorageEntry, error) {
	localPath := filepath.Join(store.root, path)

	files, err := ioutil.ReadDir(localPath)
	if err != nil {
		return nil, err
	}

	entries := make([]storage.StorageEntry, 0, len(files))

	for _, f := range files {
		entries = append(entries, storage.StorageEntry{
			Title:    f.Name(),
			Path:     filepath.Join(store.root, f.Name()),
			IsDir:    f.IsDir(),
			Updated:  f.ModTime(),
			MimeType: getMime(f.Name()),
		})
	}
	return entries, nil
}

func (store *FsStorage) IsDownloadable(path string) (bool, error) {
	localPath := filepath.Join(store.root, path)

	f, err := os.Open(localPath)
	if err != nil {
		return false, err
	}

	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return false, err
	}

	return !stat.IsDir(), nil
}

func (store *FsStorage) Download(w io.Writer, path string) error {
	localPath := filepath.Join(store.root, path)

	f, err := os.Open(localPath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(w, f)

	return err
}

func getMime(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType, ok := mimes[ext]
	if ok {
		return mimeType
	}

	mimeType = mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return mimeType
}
