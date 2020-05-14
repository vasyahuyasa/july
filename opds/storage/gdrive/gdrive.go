package gdrive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vasyahuyasa/july/opds/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	MIMETypeFolder = "application/vnd.google-apps.folder"

	tokFile = "token.json"
)

var _ storage.Storage = &GdriveStorage{}

type GdriveStorage struct {
	rootID string
	svc    *drive.Service
}

func NewStorage(root string, svc *drive.Service) (*GdriveStorage, error) {
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

			mod, err := time.Parse(time.RFC3339Nano, i.ModifiedTime)
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
	fileID := drive.filterpath(path)

	resp, err := drive.svc.Files.Get(fileID).Download()
	if err != nil {
		return fmt.Errorf("can not download %q: %w", fileID, err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("can not send file content: %w", err)
	}

	return nil
}

func (drive *GdriveStorage) filterpath(path string) string {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	if path == "" {
		return drive.rootID
	}

	return path
}

func OAuth2FromFile(filename string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %w", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %w", err)
	}

	return config, nil
}

func NewServiceFromOauth2(config *oauth2.Config) (*drive.Service, error) {
	client := getClient(config)

	srv, err := drive.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %w", err)
	}

	return srv, err
}

func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
