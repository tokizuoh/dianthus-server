package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	storageBucket := os.Getenv("FIREBASE_STORAGE_BUCKET")

	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	cred := os.Getenv("FIREBASE_CREDENTIAL_FILE_PATH")
	opt := option.WithCredentialsFile(cred)
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatal(err)
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatal(err)
	}

	fp := os.Getenv("FIREBASE_STORAGEL_FILE_PATH")
	ctx := context.Background()
	rc, err := bucket.Object(fp).NewReader(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(data))
}
