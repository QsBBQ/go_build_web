package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
  DBHost  = "127.0.0.1"
  DBPort  = ":5432"
  DBUser  = "gowebuser"
  DBPass  = "mypassword1"
  DBDbase = "goweb_practice"
  PORT    = ":8080"
)

var database *sql.DB

type Page struct {
	Title   string
	Content string
	Date    string
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageID := vars["id"]
	thisPage := Page{}
	fmt.Println(pageID)
	err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE id=$1", pageID).Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)
  fmt.Println(err)
  if err != nil {
		log.Println(err)
		log.Println("Couldn't get page!")
	}
	html := `<html><head><title>` + thisPage.Title + `</title></head><body><h1>` + thisPage.Title + `</h1><div>` + thisPage.Content + `</div></body></html>`
	fmt.Fprintln(w, html)
}

func main() {
  // db, err := sql.Open("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
	dbConn := fmt.Sprintf("postgres://%s:%s@/%s?sslmode=disable", DBUser, DBPass, DBDbase)
	fmt.Println(dbConn)
	db, err := sql.Open("postgres", dbConn)
  // fmt.Println(db)
	if err != nil {
		log.Println("Couldn't connect!")
		log.Println(err.Error)
	}
	database = db

	routes := mux.NewRouter()
	routes.HandleFunc("/page/{id:[0-9]+}", ServePage)
	http.Handle("/", routes)
	http.ListenAndServe(PORT, nil)

}
