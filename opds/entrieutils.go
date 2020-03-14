package opds

import (
	"path/filepath"

	"github.com/vasyahuyasa/july/opds/storage"
)

const (
	dirMime = "application/atom+xml;profile=opds-catalog;kind=navigation"
	dirRel  = "subsection"
	fileRel = "http://opds-spec.org/acquisition"
)

func entriesFromStorage(list []storage.StorageEntry, basePath string) []Entry {
	entries := make([]Entry, 0, len(list))

	var e Entry

	for _, se := range list {
		if se.IsDir {
			e = dirEntry(se, basePath)
		} else {
			e = fileEntry(se, basePath)
		}
		entries = append(entries, e)
	}
	return entries
}

func dirEntry(e storage.StorageEntry, basePath string) Entry {
	link := filepath.Join(basePath, e.Path)

	return Entry{
		ID:      link,
		Updated: e.Updated.UTC().Format(atomTime),
		Title:   e.Title,
		Link: []Link{
			Link{
				Href: link,
				Type: dirMime,
				Rel:  dirRel,
			},
		},
	}
}

func fileEntry(e storage.StorageEntry, basePath string) Entry {
	link := filepath.Join(basePath, e.Path)

	return Entry{
		ID:      link,
		Updated: e.Updated.UTC().Format(atomTime),
		Title:   e.Title,
		Link: []Link{
			Link{
				Href: link,
				Type: e.MimeType,
				Rel:  fileRel,
			},
		},
	}
}
