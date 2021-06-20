package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type Word struct {
	Raw    string `json:"raw"`
	Roman  string `json:"roman"`
	Vowels string `json:"vowels"`
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

func getWordsWithSameVowel(target string, data []byte) []Word {
	var result []Word
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

		w := Word{
			Raw:    _line[1],
			Roman:  _line[2],
			Vowels: vowels,
		}
		result = append(result, w)
	}
	return result
}

func validAuth(r *http.Request) bool {
	reqCliID, reqCliSec, ok := r.BasicAuth()
	if ok != true {
		return false
	}

	if err := godotenv.Load(); err != nil {
		return false
	}

	cliID := os.Getenv("BASIC_AUTH_CLIENT_ID")
	cliSec := os.Getenv("BASIC_AUTH_CLIENT_SECRET")
	return cliID == reqCliID && cliSec == reqCliSec
}

func main() {
	h := func(w http.ResponseWriter, r *http.Request) {

		if validAuth(r) != true {
			http.Error(w, "Unauthorized", 401)
			return
		}

		q := r.URL.Query().Get("target")
		if q == "" {
			http.Error(w, "Bad Request", 400)
			return
		}

		data, err := fetchCSV()
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		target := extractVowels(q)

		words := getWordsWithSameVowel(target, data)

		res, err := json.Marshal(words)
		w.Header().Set("Content-Type", "application/json")
		if _, err = w.Write(res); err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
	}

	http.HandleFunc("/v1/roman", h)
	http.ListenAndServe(":8080", nil)
}
