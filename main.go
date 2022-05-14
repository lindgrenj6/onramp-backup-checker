package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"
)

func main() {
	if os.Getenv("DROPBOX_TOKEN") == "" {
		log.Fatalf("need DROPBOX_TOKEN in order to run")
	}
	log.SetFlags(log.Lshortfile)

	config := dropbox.Config{Token: os.Getenv("DROPBOX_TOKEN")}
	f := files.New(config)

	res := must(f.ListFolder(files.NewListFolderArg("")))
	names := make([]string, 0)

	for _, f := range res.Entries {
		switch f.(type) {
		case *files.FileMetadata:
			file := f.(*files.FileMetadata)
			if strings.HasPrefix(file.Name, "onramp_production") {
				names = append(names, file.Name)
			}
		default:
			log.Printf("skipping %v", reflect.TypeOf(f))
		}
	}

	sort.Strings(names)
	now := time.Now()
	filename := fmt.Sprintf("onramp_production_%v%02d%02d.sql.gz", now.Year(), now.Month(), now.Day())

	if names[len(names)-1] != filename {
		log.Fatalf("today's backup is missing!")
	}
}

func must[T any](thing T, err error) T {
	if err != nil {
		panic(err)
	}
	return thing
}
