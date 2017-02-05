package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"html/template"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	PORT    = ":8080"
	DBHost  = "127.0.0.1"
	DBPort  = ":5432"
	DBUser  = "postgres"
	DBPass  = "postgres"
	DBDbase = "golang"
)

var database *sql.DB

type Page struct {
	GUID 				string
	Title   		string
	RawContent 	string
	Content 		template.HTML
	Date    		string
}

func (p Page) TruncantedText() template.HTML {
	chars := 0
	for i, _ := range p.Content {
		chars++
		if chars > 5 {
			return p.Content[:i] + `...`
		}
	}
	return p.Content
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	err := database.QueryRow(`SELECT page_title, 
														page_content FROM 
														pages WHERE page_guid=$1`, pageGUID).Scan(&thisPage.Title, &thisPage.RawContent)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page: " + pageGUID)
	}

	thisPage.Content = template.HTML(thisPage.RawContent)

	t, _ := template.ParseFiles("templates/blog.html")
  t.Execute(w, thisPage)
}

func RedirIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 301)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	var Pages = []Page{}
	pages, err := database.Query(`SELECT page_title,
																page_content, page_guid FROM
																pages ORDER BY $1 DESC`, "page_date")
	if err != nil {
		fmt.Fprintln(w, err.Error)
	}

	defer pages.Close()
	for pages.Next() {
		thisPage := Page{}
		pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.GUID)
		thisPage.Content = template.HTML(thisPage.RawContent)
		Pages = append(Pages, thisPage)
	}
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, Pages)
}

func main() {
	dbConn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DBUser, DBPass,
		DBHost, DBDbase)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Println("Couldn't connect!")
		log.Println(err.Error)
	}
	database = db

	routes := mux.NewRouter()
	routes.HandleFunc("/pages/{guid:[0-9a-zA\\-]+}", ServePage)
	routes.HandleFunc("/", RedirIndex)
	routes.HandleFunc("/home", ServeIndex)
	http.Handle("/", routes)
	http.ListenAndServe(PORT, nil)
}
