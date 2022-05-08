package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Post struct {
	UserId int
	Id     int
	Title  string
	Body   string
}

func main() {
	baseUrl := "https://jsonplaceholder.typicode.com/posts/"

	for i := 1; i <= 5; i++ {
		go performRequest(baseUrl + strconv.Itoa(i))
	}

	fmt.Scanln()
}

func performRequest(url string) {
	response, error := http.Get(url)

	checkError(error)

	newPost := parseJsonToStruct(response)

	saveToFile(newPost, url[len(url)-1:])

	defer response.Body.Close()
}

func saveToFile(post Post, index string) {
	structBytes := new(bytes.Buffer)

	json.NewEncoder(structBytes).Encode(post)

	ioutil.WriteFile("./storage/posts/"+index+".txt", structBytes.Bytes(), 0666)
}

func parseJsonToStruct(response *http.Response) Post {
	jsonDataFromHttp, ioError := ioutil.ReadAll(response.Body)

	checkError(ioError)

	var post Post

	jsonError := json.Unmarshal(jsonDataFromHttp, &post)

	checkError(jsonError)

	return post
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
