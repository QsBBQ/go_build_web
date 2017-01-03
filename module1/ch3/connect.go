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
  pageID := vars["id"]
  thisPage := Page{}
  fmt.Println(pageID)
  err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE id=$1", pageID).Scan(&thisPage.Title, &thisPage.Content, &thisPage.Date)
  // err := database.QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE id=?", pageID).Scan(&thisPage.Title,
                          //  &thisPage.Content, &thisPage.Date)
  if err != nil {
    log.Println("Couldn't get page: "+pageID)
    log.Println(err)
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
  routes.HandleFunc("/page/{id:[0-9]+}", ServePage)
  http.Handle("/", routes)
  http.ListenAndServe(PORT, nil)
}

// func main() {
//   rtr := mux.NewRouter()
//   // rtr.HandleFunc("/pages/{id:[0-9]+}", pageHandler)
//   http.Handle("/",rtr)
//   http.ListenAndServe(PORT, nil)
// }
