package main

import (
	"context"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	_ "github.com/lib/pq"
	"github.com/orov.io/GoBackbone/service"
	"github.com/sirupsen/logrus"
)

func init() {
	profile := os.Getenv("ENV")

	if profile == "local" {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			FullTimestamp:   true,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		ctx := context.Background()
		bucket := os.Getenv("GCLOUD_STORAGE_BUCKET")

		storageClient, err := storage.NewClient(ctx)
		if err != nil {
			logrus.Fatalf("Can't get storage conection: %v", err)
		}
		sw, err := storageClient.Bucket(bucket).Object("firebase_credential.json").NewReader(ctx)
		if err != nil {
			logrus.Fatalf("Can't stablish reader connection: %v", err)
		}
		buff, err := ioutil.ReadAll(sw)
		if err != nil {
			logrus.Fatalf("Can't read credentials: %v", err)
		}
		err = ioutil.WriteFile("firebase_credential.json", buff, 0644)
		if err != nil {
			logrus.Fatalf("Can't write firebase_credential.json: %v", err)
		}
	}
}

func main() {
	app := service.App{}
	app.Initialize()

	// Start and run the server
	app.Run()
}
