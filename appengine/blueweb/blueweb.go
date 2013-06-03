package blueweb 

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/get10", get10)
} 

type Get10Options struct {
  NoAttack bool
  NoAttackWithoutDefense bool
  SetsAvailable []int 
  ForbiddenCards []int
}

func readStructFromJSONRequest(w http.ResponseWriter, r *http.Request, readInto interface{}) error {
    // TODO: we should probably check that this is a POST and that it's a request of type text/json
  jsonRaw,_ := ioutil.ReadAll(r.Body)
  err := json.Unmarshal(jsonRaw, readInto)
  if err != nil {
    serveError(w, err)
  }
  return err
}

// this is where all the magic happens for generating the 10 cards
func generateCards(options *Get10Options) []int {
  return []int{1, 2, 3, 7, 9, 4, 0, 8, 6, 5}
}

func generateGet10Response(cards []int, w http.ResponseWriter) {
  b,_ := json.Marshal(cards)
  fmt.Fprint(w, string(b))
}

func get10(w http.ResponseWriter, r *http.Request) {
  options := Get10Options{}
  optionsErr := readStructFromJSONRequest(w, r, &options)
  if optionsErr != nil {
    return
  }

  cards := generateCards(&options) 

  generateGet10Response(cards, w)
}

func root(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Welcome to Just the Blue Cards")
}

func serveError(w http.ResponseWriter, err error) {
    w.WriteHeader(500)
    fmt.Fprintf(w, 
                "whoopsies! Could not understand that json struct:\n'%v'",
                err)
}


