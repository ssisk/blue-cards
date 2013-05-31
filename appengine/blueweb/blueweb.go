package blueweb 

import (
    "encoding/json"
    "fmt"
    "net/http"
)

func init() {
    http.HandleFunc("/", root)
}

type Message struct {
    Name string
    Body string
    Time int64
}

func root(w http.ResponseWriter, r *http.Request) {
  val := r.FormValue("val")

  m := Message{val, "Hello", 1294706395881547000}

  b,err := json.Marshal(m)
  err = err

  fmt.Fprint(w, string(b))

}