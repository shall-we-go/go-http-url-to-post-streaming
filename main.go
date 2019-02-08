package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

func main() {
	token := "ZjRhY2EwNjg4NmVhY2ViYTViMWY2Y2JlYmFlOTQzZjhkZGJmYjIzYTFhY2FiMTU3NGU2MGNlNTcwMjJhYjAxMQ"
	fileUrl := "https://media-dev-eu-2.mirakl.net/SOURCE/48bacd161d5f49bcac123034f27b844f"
	apiURL := "http://dev-pim.central.co.th/api/rest/v1/media-files"
	log.Println("Downloading " + fileUrl)

	resp, err := http.Get(fileUrl)
	if err != nil {
		log.Fatal(err)
		return
	}

	fileReader := bufio.NewReader(resp.Body)
	formValues := map[string]io.Reader{
		"product": strings.NewReader(`{"identifier":"MKT2000006", "attribute":"example_image", "scope": null,"locale":null}`),
		"file":    fileReader,
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range formValues {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		if key == "file" {
			if fw, err = w.CreateFormFile(key, "file"); err != nil {
				log.Fatal(err)
			}
		} else {
			if fw, err = w.CreateFormField(key); err != nil {
				log.Fatal(err)
			}
		}

		if _, err = io.Copy(fw, r); err != nil {
			log.Fatal(err)
		}

	}
	w.Close()

	req, err := http.NewRequest("POST", apiURL, &b)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", body)
	}

}
