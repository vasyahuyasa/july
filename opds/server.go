package opds

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const startURL = "/opds"
const title = "July: home opds catalog"
const atomTime = "2006-01-02T15:04:05Z"

type Server struct {
	// FileRoot is root catalog of file storage
	// by default it is current directory
	FileRoot string
}

func (s *Server) opdsHandler(w http.ResponseWriter, r *http.Request) {
	currentURL := "http://" + r.Host + r.URL.Path
	start := "http://" + r.Host + startURL
	path := r.URL.Path[len(startURL):]

	// default file root is current directory
	if s.FileRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Println("os.Getwd():", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "os.Getwd(): %s", err)
			return
		}
		s.FileRoot = cwd
	}

	root := s.FileRoot + path

	// check for listing or download
	info, err := os.Stat(root)
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("file %q not exists\n", root)
		fmt.Fprintf(w, "file %q not exists", path)
		return
	}

	if !info.IsDir() {
		log.Printf("Download %q", root)
		w.Header().Add("Content-Disposition", "Attachment")
		http.ServeFile(w, r, root)
		return
	}

	log.Printf("Linsting for %q\n", root)

	entries, err := makeEntries(root, r.URL.Path)
	if err != nil {
		log.Printf("makeEntries(%q, %q): %v\n", root, currentURL, err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "makeEntries(%q, %q): %v", root, currentURL, err)
		return
	}

	f := &Feed{
		ID:      currentURL,
		Title:   title,
		Xmlns:   "http://www.w3.org/2005/Atom",
		Updated: time.Now().UTC().Format(atomTime),
		Link: []*Link{
			&Link{
				Href: start,
				Type: "application/atom+xml;profile=opds-catalog;kind=navigation",
				Rel:  "start",
			},
			&Link{
				Href: r.URL.Path,
				Type: "application/atom+xml;profile=opds-catalog;kind=navigation",
				Rel:  "self",
			},
		},
		Entry: entries,
	}

	enc := xml.NewEncoder(w)

	w.Header().Set("Content-Type", "application/xml")
	fmt.Fprint(w, xml.Header)
	err = enc.Encode(f)
	if err != nil {
		log.Println("enc.Encode(v): ", err)
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
