package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type uploadResponse struct {
	File   string `json:"file"`
	Upload string `json:"upload"`
}

func main() {
	for _, file := range os.Args {
		uri, err := uploadFile(context.Background(), file)
		if err != nil {
			log.Printf("error: could not upload file %q: %+v", file, err)
			os.Exit(1)
		}

		fmt.Println(uri)
	}
}

func uploadFile(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %w", err)
	}

	var b []byte
	if _, err := file.Read(b); err != nil {
		return "", fmt.Errorf("could not read file: %w", err)
	}
	mimeType := http.DetectContentType(b)
	buf := bytes.NewBuffer(b)

	resp, err := http.Post("https://graphql.natwelch.com/photo/new", mimeType, buf)
	if err != nil {
		return "", fmt.Errorf("could not upload file: %w", err)
	}

	var upload uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&upload); err != nil {
		return "", fmt.Errorf("could not decode response: %w", err)
	}

	return upload.File, nil
}
