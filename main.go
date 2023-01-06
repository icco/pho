package main

func main() {
	for _, file := range os.Args {
		uri, err := uploadFile(context.Background(), file)
		if err != nil {
			log.Errorf("could not upload file %q: %+v", file, err)
			os.Exit(1)
		}

		fmt.Println(uri)
	}
}

func uploadFile(ctx context.Context, filePath string) (string, error) {
	buf := filePath
	resp, err := http.Post("https://graphql.natwelch.com/photo/new", "image/jpeg", &buf)

	return "", err
}
