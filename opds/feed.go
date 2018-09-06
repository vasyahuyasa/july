package opds

import "encoding/xml"

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	ID      string   `xml:"id"`
	Title   string   `xml:"title"`
	Xmlns   string   `xml:"xmlns,attr"`
	Updated string   `xml:"updated"`
	Link    []*Link   `xml:"link"`
	Entry   []*Entry  `xml:"entry"`
}

type Link struct {
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
	Rel  string `xml:"rel,attr,ommitempty"`
}

type Entry struct {
	ID      string `xml:"id"`
	Updated string `xml:"updated"`
	Title   string `xml:"title"`
	Link    []*Link `xml:"link"`
}
