package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

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
	Title   string
	Content string
	Date    string
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageGUID := vars["guid"]
	thisPage := Page{}
	err := database.QueryRow(`SELECT page_title, 
														page_content FROM 
														pages WHERE page_guid=$1`, pageGUID).Scan(&thisPage.Title, &thisPage.Content)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page: " + pageGUID)
	}

	html := `<html><head><title>` + thisPage.Title +
		`</title></head><body><h1>` + thisPage.Title + `</h1><div>` +
		thisPage.Content + `</div></body></html>`

	fmt.Fprintln(w, html)
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
	http.Handle("/", routes)
	http.ListenAndServe(PORT, nil)
}
