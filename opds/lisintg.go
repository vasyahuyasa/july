package opds

import (
	"io/ioutil"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var mimes = map[string]string{
	".epub": "application/epub+zip",
	".fb2":  "application/fb2+zip",
	".mobi": "application/x-mobipocket-ebook",
}

func makeEntries(root, rootURL string) ([]*Entry, error) {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	entries := make([]*Entry, 0, len(files))

	for _, f := range files {
		var e *Entry
		if f.IsDir() {
			e = dirEntry(f, rootURL+"/"+url.PathEscape(f.Name()))
		} else {
			e = fileEntry(f, rootURL+"/"+url.PathEscape(f.Name()))
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func dirEntry(info os.FileInfo, URL string) *Entry {
	return &Entry{
		ID:      URL,
		Updated: info.ModTime().UTC().Format(atomTime),
		Title:   info.Name(),
		Link: []*Link{
			&Link{
				Href: URL,
				Type: "application/atom+xml;profile=opds-catalog;kind=navigation",
				Rel:  "subsection",
			},
		},
	}
}

func fileEntry(info os.FileInfo, URL string) *Entry {
	return &Entry{
		ID:      URL,
		Updated: info.ModTime().UTC().Format(atomTime),
		Title:   info.Name(),
		Link: []*Link{
			&Link{
				Href: URL,
				Type: getMime(info.Name()),
				Rel:  "http://opds-spec.org/acquisition",
			},
		},
	}
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
