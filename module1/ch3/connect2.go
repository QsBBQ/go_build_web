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
  Title string
  Content string
  Date string
}

func ServePage(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  pageGUID := vars["guid"]
  thisPage := Page{}
  fmt.Println(pageGUID)
  err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=$1", pageGUID).Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)

  if err != nil {
    log.Println("Couldn't get page: "+pageGUID)
    log.Println(err)
    http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
		return
    }
  html := `<html><head><title>` + thisPage.Title + `</title></head><body><h1>` + thisPage.Title + `</h1><div>` + thisPage.Content + `</div></body></html>`
  fmt.Fprintln(w, html)
}
func main() {
  // db, err := sql.Open("mysql", "<username>:<pw>@tcp(<HOST>:<port>)/<dbname>")
  dbConn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DBUser,
                        DBPass, DBHost, DBDbase)
  fmt.Println("connecting to db", DBDbase, DBHost)
  db, err := sql.Open("postgres", dbConn)
  if err != nil {
    log.Println("Couldn't connect! " + DBDbase)
    log.Println(err.Error)
  }
  database = db
  routes := mux.NewRouter()
  routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
  http.Handle("/", routes)
  http.ListenAndServe(PORT, nil)
}
