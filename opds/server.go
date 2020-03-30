package opds

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vasyahuyasa/july/opds/storage"
)

const (
	startURL = "/opds"
	title    = "July: home opds catalog"
	atomTime = "2006-01-02T15:04:05Z"
)

type Server struct {
	store storage.Storage
}

func NewServer(store storage.Storage) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) opdsHandler(w http.ResponseWriter, r *http.Request) {
	currentURL := "http://" + r.Host + r.URL.Path
	start := "http://" + r.Host + startURL
	path := r.URL.Path[len(startURL):]

	log.Printf("Path %q", path)

	canDownload := false

	// root always must be folder
	if path != "" {
		var err error
		if canDownload, err = s.store.IsDownloadable(path); err != nil {
			log.Printf("can not check if path %q is downloadable: %#v", path, err)
			code := http.StatusInternalServerError
			if os.IsNotExist(err) {
				code = http.StatusNotFound
			}
			http.Error(w, fmt.Sprintf("can not check if path %q is downloadable: %v", path, err), code)
			return
		}
	}

	// if we can't download so it should be folder, lets list it
	if !canDownload {
		list, err := s.store.List(path)
		if err != nil {
			log.Printf("can not list files in %q: %v", path, err)
			http.Error(w, fmt.Sprintf("can not list file is %q: %v", path, err), http.StatusInternalServerError)
			return
		}

		f := &Feed{
			ID:      currentURL,
			Title:   title,
			Xmlns:   "http://www.w3.org/2005/Atom",
			Updated: time.Now().UTC().Format(atomTime),
			Link: []Link{
				Link{
					Href: start,
					Type: dirMime,
					Rel:  "start",
				},
				Link{
					Href: r.URL.Path,
					Type: dirMime,
					Rel:  "self",
				},
			},
			Entry: entriesFromStorage(list, start),
		}

		enc := xml.NewEncoder(w)

		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, xml.Header)
		err = enc.Encode(f)
		if err != nil {
			log.Printf("can not encode xml feed: %v", err)
			if err != io.ErrClosedPipe {
				http.Error(w, fmt.Sprintf("can not encode xml feed: %v", err), http.StatusInternalServerError)
			}
		}
		return
	}

	// download file
	err := s.store.Download(w, path)
	if err != nil {
		log.Printf("can not download %q: %v", path, err)
		if err != io.ErrClosedPipe {
			http.Error(w, fmt.Sprintf("can not download %q: %v", path, err), http.StatusInternalServerError)
		}
	}
}

func (s *Server) Run(addr string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request %q: %q", r.RemoteAddr, r.URL.String())

		if len(r.URL.Path) >= len(startURL) && r.URL.Path[:len(startURL)] == startURL {
			s.opdsHandler(w, r)
			return
		}
		fmt.Fprintf(w, `
<html>
	<head>
		<title>Welcome to July opds catalog</title>
	</head>
	<h1>Welcome to July opds catalog</h1>
	<a href="/opds">Opds catalog</a>
</html>`)
	})
	return http.ListenAndServe(addr, nil)
}
