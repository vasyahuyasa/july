package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/vasyahuyasa/july/opds"
	"github.com/vasyahuyasa/july/opds/storage"
	"github.com/vasyahuyasa/july/opds/storage/gdrive"
	"github.com/vasyahuyasa/july/opds/storage/local"
)

func main() {
	root := flag.String("d", "./", "Root storage directory")
	port := flag.Int("p", 80, "Service port")
	host := flag.String("i", "0.0.0.0", "Service network interface")
	driver := flag.String("drv", "local", "Storage driver (can be local, gdrive, yadisk)")

	//gdriveCredFile := flag.String("gcred", "credentials.json", "Path to file with secret for google drive driver")
	//gdriveTokenFile := flag.String("gtoken", "token.json", "Path to file with token for google drive driver")
	key := flag.String("k", "", "API key for cloud storage providers (gdrive, yadisk)")

	flag.Parse()

	var store storage.Storage

	switch *driver {
	case "local":
		store = local.NewFsStorage(*root)
	case "gdrive":
		var err error

		store, err = gdrive.NewStorage(*root, *key)
		if err != nil {
			log.Fatalf("Can not initialize gdrive: %v", err)
		}
	case "yadisk":
		panic(*driver + " not implemented yet")
	default:
		panic("unknown driver " + *driver)
	}

	srv := opds.NewServer(store)

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Run server: ", "http://"+addr+"/opds")
	err := srv.Run(addr)
	log.Fatal(err)
}
