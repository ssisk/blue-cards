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

func get10ReadOptions (w http.ResponseWriter, r *http.Request) (Get10Options, error) {
  // TODO: we should probably check that this is a POST and that it's a request of type text/json
  options := Get10Options{}
  optionsRaw,_ := ioutil.ReadAll(r.Body)
  err := json.Unmarshal(optionsRaw, &options)
  // just a little debug code - feel free to remove eventually
  //optionsString := `{"NoAttack": true,"NoAttackWithoutDefense": false}`
  //err := json.Unmarshal([]byte(optionsString), &options)
  if err != nil {
    w.WriteHeader(500)
    fmt.Fprint(w, "Those options were not understood\n")
    fmt.Fprintf(w, "optionsRaw:'%v'\n", string(optionsRaw))
  }
  return options, err
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
  options, optionsErr := get10ReadOptions(w, r)
  if optionsErr != nil {
    return
  }

  cards := generateCards(&options) 

  generateGet10Response(cards, w)
}


func root(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Welcome to Just the Blue Cards")
}