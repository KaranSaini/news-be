package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/alexsasharegan/dotenv"
	"github.com/gorilla/mux"
)

type NewsResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Source struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

var dataToShare []Article

func news(w http.ResponseWriter, r *http.Request) {
	category := mux.Vars(r)
	getNewsArticles(category["category"])
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
}

// function that should get articles from NewsAPI and save them
func getNewsArticles(category string) {
	newsToken := os.Getenv("NEWS_TK")
	url := fmt.Sprintf("http://newsapi.org/v2/top-headlines?category=%v&country=us&pageSize=100&apiKey=%v", category, newsToken)
	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("api key limit has probably been hit")
		log.Fatalln("error", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("error", err)
	}
	// this works and I get a type of NewsResponse from data ... now I want to loop through
	var data NewsResponse
	errU := json.Unmarshal(body, &data)
	if errU != nil {
		log.Fatalln(errU)
	}
	dataToShare = data.Articles
}

func main() {
	err := dotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	//Want to Schedule News Collection Twice a Day

	r := mux.NewRouter()
	r.HandleFunc("/news/{category}", news)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
