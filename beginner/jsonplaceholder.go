package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const url string = "https://jsonplaceholder.typicode.com/posts"

func main() {
	response, error := http.Get(url)

	if error != nil {
		fmt.Println("Error: ", error)
		os.Exit(1)
	}

	io.Copy(os.Stdout, response.Body)
}
