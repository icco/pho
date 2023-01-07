package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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
	r, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}
	defer r.Close()

	client := &http.Client{
		Transport: &AddHeaderTransport{T: http.DefaultTransport, Key: apikey},
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	key := "file"
	fw, err := w.CreateFormFile(key, r.Name())
	if err != nil {
		return "", fmt.Errorf("could not create form file: %w", err)
	}

	if _, err := io.Copy(fw, r); err != nil {
		return "", fmt.Errorf("could not copy file: %w", err)
	}

	if err := w.Close(); err != nil {
		return "", fmt.Errorf("could not close writer: %w", err)
	}

	resp, err := client.Post("https://graphql.natwelch.com/photo/new", w.FormDataContentType(), &b)
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
