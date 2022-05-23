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
	"github.com/slack-go/slack"
)

func main() {
	if os.Getenv("DROPBOX_TOKEN") == "" {
		log.Fatalf("need DROPBOX_TOKEN in order to run")
	}
	if os.Getenv("SLACK_WEBHOOK_URL") == "" {
		log.Fatalf("need SLACK_WEBHOOK_URL in order to run")
	}
	log.SetFlags(log.Lshortfile)

	config := dropbox.Config{Token: os.Getenv("DROPBOX_TOKEN")}
	f := files.New(config)

	res, err := f.ListFolder(files.NewListFolderArg(""))
	if err != nil {
		log.Fatal(err)
	}
	names := make([]string, 0)

	for _, f := range res.Entries {
		switch file := f.(type) {
		case *files.FileMetadata:
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
		err := slack.PostWebhook(os.Getenv("SLACK_WEBHOOK_URL"), &slack.WebhookMessage{
			Username:  "Backup WatchDog",
			IconEmoji: ":guide_dog:",
			Channel:   "#alerts",
			Text:      "Today's backup is missing from dropbox - best check it out.",
		})

		if err != nil {
			log.Fatal(err)
		}
	}
}
