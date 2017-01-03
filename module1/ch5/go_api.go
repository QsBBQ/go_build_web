package main

import (
  "database/sql"
  "fmt"
  _ "github.com/lib/pq"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "html/template"
  "encoding/json"
  _ "crypto/tls"
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
  RawContent string
  Content template.HTML
  Date string
  GUID string
}

func ServePage(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  pageGUID := vars["guid"]
  thisPage := Page{}
  fmt.Println(pageGUID)
  err := database.QueryRow("SELECT page_title,page_content,page_date FROM pages WHERE page_guid=$1", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
  thisPage.Content = template.HTML(thisPage.RawContent)
  if err != nil {
    log.Println("Couldn't get page: "+pageGUID)
    log.Println(err)
    http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page!")
		return
    }
  t, _ := template.ParseFiles("templates/blog.html")
  t.Execute(w, thisPage)
}

func RedirIndex(w http.ResponseWriter, r *http.Request) {
  http.Redirect(w, r, "/home", 301)
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
  var Pages = []Page{}
  pages, err := database.Query("SELECT page_title, page_content, page_date, page_guid FROM pages ORDER BY $1 DESC", "page_date")
  if err != nil {
    fmt.Fprintln(w, err.Error)
    }
    defer pages.Close()
  for pages.Next() {
    thisPage := Page{}
    pages.Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date, &thisPage.GUID)
    thisPage.Content = template.HTML(thisPage.RawContent)
    Pages = append(Pages, thisPage)
    }
  t, _ := template.ParseFiles("templates/index.html")
  t.Execute(w, Pages)
}

func (p Page) TruncatedText() string {
  chars := 0
  for i, _ := range p.RawContent {
    chars++
    if chars > 150 {
      return p.RawContent[:i] + `...`
    }
  }
  return p.RawContent
}

func APIPage(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  pageGUID := vars["guid"]
  thisPage := Page{}
  fmt.Println(pageGUID)
  err := database.QueryRow("SELECT page_title, page_content, page_date FROM pages WHERE page_guid=$1", pageGUID).Scan(&thisPage.Title, &thisPage.RawContent, &thisPage.Date)
  thisPage.Content = template.HTML(thisPage.RawContent)
  if err != nil {
    fmt.Println("Trouble in paradise")
    http.Error(w, http.StatusText(404), http.StatusNotFound)
    log.Println(err)
    return
  }
  APIOutput, err := json.Marshal(thisPage)
  fmt.Println(APIOutput)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintln(w, thisPage)
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
  routes.HandleFunc("/api/pages", APIPage).Methods("GET").Schemes("https")
  routes.HandleFunc("/api/page/{guid:[0-9a-zA\\-]+}", APIPage).Methods("GET")
  .Methods("GET").Schemes("https")
  routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)
  routes.HandleFunc("/", RedirIndex)
  routes.HandleFunc("/home", ServeIndex)
  http.Handle("/", routes)
  http.ListenAndServe(PORT, nil)

  // This does not work with ssl?
  // certificates, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
  // if err != nil {
  //   log.Println(err)
  //   return
  // }
  // tlsConf := tls.Config{Certificates: []tls.Certificate{certificates}}
  // tls.Listen("tcp", PORT, &tlsConf)

  // This works with ssl
  // http.ListenAndServeTLS(PORT, "cert.pem", "key.pem", nil)
}
