package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type word struct {
	raw    string
	roman  string
	vowels string
}

func fetchCSV() ([]byte, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	storageBucket := os.Getenv("FIREBASE_STORAGE_BUCKET")

	config := &firebase.Config{
		StorageBucket: storageBucket,
	}
	cred := os.Getenv("FIREBASE_CREDENTIAL_FILE_PATH")
	opt := option.WithCredentialsFile(cred)
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, err
	}

	fp := os.Getenv("FIREBASE_STORAGEL_FILE_PATH")
	ctx := context.Background()
	rc, err := bucket.Object(fp).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func extractVowels(from string) string {
	vowel := "aiuoe"

	result := ""
	shouldSkip := false

	for i, sr := range from {
		if shouldSkip {
			shouldSkip = false
			continue
		}

		s := string([]rune{sr})

		if strings.Contains(vowel, s) {
			result += s
			continue
		}

		if s == "n" {
			if i+1 < len(from) {
				ns := from[i+1 : i+2]
				if strings.Contains(vowel, ns) {
					shouldSkip = true
					result += ns
				} else {
					result += s
				}
			} else {
				result += s
			}
			continue
		}

		if i+1 < len(from) {
			ns := from[i+1 : i+2]
			if s == ns {
				result += "x"
			}
		} else {
			result += s
		}
	}

	return result
}

func getWordsWithSameVowel(target string, data []byte) []word {
	var result []word
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		_line := strings.Split(line, ",")

		if len(_line) != 4 {
			continue
		}

		vowels := strings.TrimSpace(_line[3])

		if vowels != target {
			continue
		}

		w := word{
			raw:    _line[1],
			roman:  _line[2],
			vowels: _line[3],
		}
		result = append(result, w)
	}
	return result
}

func main() {
	data, err := fetchCSV()
	if err != nil {
		log.Fatal(err)
	}

	target := extractVowels("kien")

	words := getWordsWithSameVowel(target, data)

	for _, r := range words {
		log.Println(r.raw, r.roman, r.vowels)
	}

	log.Println("FINISH!")
}
