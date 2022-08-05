package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// post()
	get()
}
func post() {
	var URL string
	fmt.Scanf("%s\n", &URL)

	req, err := http.NewRequest("POST", "http://localhost:8080", bytes.NewBufferString(URL))
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req.Header.Set("URL", URL)
	resp, err := client.Do(req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(body))
	fmt.Println(resp.StatusCode)
	fmt.Println(req.Header)
}

func get() {
	client := &http.Client{}

	req1, err := http.NewRequest("GET", "http://localhost:8080/query?id=XVlBzgbaiC", nil)
	q := req1.URL.Query().Get("id")
	fmt.Println(q)
	fmt.Println(q)

	req1.Header.Set("Location", req1.URL.Path)

	if err != nil {
		log.Fatal(err)
	}
	resp1, err := client.Do(req1)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!11")
	fmt.Println(req1.URL)
	fmt.Println(resp1.Header)
	fmt.Println(resp1.StatusCode)
}
