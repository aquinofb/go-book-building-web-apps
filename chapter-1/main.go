package main

import (
	"fmt"
	"net/http"
	"time"
)

const (
	PORT = ":8080"
)

func ServeDynamic(w http.ResponseWriter, r *http.Request) {
	response := "The time is now " + time.Now().String()
	fmt.Fprintln(w, response)
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static.html")
}

func main() {
	http.HandleFunc("/static", ServeStatic)
	http.HandleFunc("/", ServeDynamic)
	http.ListenAndServe(PORT, nil)
}
