package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type Post struct {
	UserId int
	Id     int
	Title  string
	Body   string
}

type Comment struct {
	PostId int
	Id     int
	Name   string
	Email  string
	Body   string
}

const baseUrl string = "https://jsonplaceholder.typicode.com/"

func main() {
	connStr := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	checkError(err)
	defer db.Close()

	response, httpErr := http.Get(baseUrl + "posts?userId=7")
	checkError(httpErr)

	posts := parsePostJsonToStructArray(response)

	c := make(chan int)

	for _, post := range posts {
		go writePostToDatabase(db, post, c)
		go getCommentsByPost(db, c)
	}

	fmt.Scanln()
}

func getCommentsByPost(db *sql.DB, c chan int) {
	postId := strconv.Itoa(<-c)

	response, error := http.Get(baseUrl + "comments?postId=" + postId)
	checkError(error)

	fmt.Println("Getting a comments...", postId)

	comments := parseCommentJsonToStructArray(response)

	for _, comment := range comments {
		go writeCommentToDatabase(db, comment)
	}
}

func writePostToDatabase(db *sql.DB, post Post, c chan int) {
	_, err := db.Exec("insert into posts (id, user_id, title, body) values ($1, $2, $3, $4)",
		post.Id, post.UserId, post.Title, post.Body)

	checkError(err)

	fmt.Println("Writing a post to database...", post.Id)

	c <- post.Id
}

func writeCommentToDatabase(db *sql.DB, comment Comment) {
	_, err := db.Exec("insert into comments (id, post_id, name, body, email) values ($1, $2, $3, $4, $5)",
		comment.Id, comment.PostId, comment.Name, comment.Body, comment.Email)

	checkError(err)

	fmt.Println("Writing a comment to database...", comment.Id)
}

func parsePostJsonToStructArray(response *http.Response) []Post {
	jsonDataFromHttp := getByteDataFromHttpRes(response)

	var posts []Post

	jsonError := json.Unmarshal(jsonDataFromHttp, &posts)

	checkError(jsonError)

	return posts
}

func parseCommentJsonToStructArray(response *http.Response) []Comment {
	jsonDataFromHttp := getByteDataFromHttpRes(response)

	var comments []Comment

	jsonError := json.Unmarshal(jsonDataFromHttp, &comments)

	checkError(jsonError)

	return comments
}

func getByteDataFromHttpRes(response *http.Response) []byte {
	jsonDataFromHttp, ioError := ioutil.ReadAll(response.Body)

	checkError(ioError)

	return jsonDataFromHttp
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
