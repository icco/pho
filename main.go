package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

type uploadResponse struct {
	File   string `json:"file"`
	Upload string `json:"upload"`
}

// AddHeaderTransport is a http transport for adding auth headers to a request.
type AddHeaderTransport struct {
	T   http.RoundTripper
	Key string
}

// RoundTrip actually adds the headers.
func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if adt.Key == "" {
		return nil, fmt.Errorf("no key provided")
	}

	req.Header.Add("X-API-AUTH", adt.Key)
	req.Header.Add("User-Agent", "etu/1.0")

	return adt.T.RoundTrip(req)
}

func main() {
	for i, file := range os.Args {
		if i == 0 {
			continue
		}

		uri, err := uploadFile(os.Getenv("GQL_TOKEN"), file)
		if err != nil {
			log.Printf("error: could not upload file %q: %+v", file, err)
			os.Exit(1)
		}

		fmt.Println(uri)
	}
}

func uploadFile(apikey, filePath string) (string, error) {
	log.Printf("attempting upload of %q", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}

	var b []byte
	if _, err := file.Read(b); err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}

	mimeType, err := mimetype.DetectFile(filePath)
	if err != nil {
		return "", fmt.Errorf("could not detect mime type: %w", err)
	}
	log.Printf("detected mime type %q", mimeType.String())

	buf := bytes.NewBuffer(b)

	client := &http.Client{
		Transport: &AddHeaderTransport{T: http.DefaultTransport, Key: apikey},
	}

	resp, err := client.Post("https://graphql.natwelch.com/photo/new", mimeType.String(), buf)
	if err != nil {
		return "", fmt.Errorf("could not upload file: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}

	var upload uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&upload); err != nil {
		return "", fmt.Errorf("could not decode response: %w", err)
	}

	return upload.File, nil
}
