package blueweb

import (
  "appengine"
  "appengine/datastore"
  "blueshared"
  "encoding/json"
  "fmt"
  "generate"
  "io/ioutil"
  "net/http"
)

const RatingsKindName = "SetRatings"

func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/get10", get10)
  http.HandleFunc("/judge", judge)
  http.HandleFunc("/analysis", analysis)
  http.HandleFunc("/matchingCard", matchingCard)
}

func readStructFromJSONRequest(w http.ResponseWriter, r *http.Request, readInto interface{}) error {
  // TODO: we should probably check that this is a POST and that it's a request of type text/json
  jsonRaw, _ := ioutil.ReadAll(r.Body)
  return json.Unmarshal(jsonRaw, readInto)
} 
 


func generateGet10Response(cards []int, w http.ResponseWriter) {
  b, _ := json.Marshal(cards)
  fmt.Fprint(w, string(b))
}

func get10(w http.ResponseWriter, r *http.Request) {
  var err error

  options := blueshared.Get10Options{}
  if err = readStructFromJSONRequest(w, r, &options); err != nil {
    serveError(w, err)
    return
  }

  cards := generate.CardSet(&options)

  generateGet10Response(cards, w)
}

func root(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "Welcome to Just the Blue Cards")
}

func serveError(w http.ResponseWriter, err error) {
  w.WriteHeader(500)
  fmt.Fprintf(w,
    "whoopsies! There was an error:\n'%v'",
    err)
}

type SetRating struct {
  /* 
    Obviously, these are going to change. This is just a first pass.

    If we can, we might just make this a map that gets some light validation and
    written straight into the db. That way, we can add new survey questions on the
    client, and the analysis  code can read it without us having to deploy the
    web layer.

  */

  /* We should talk about how to store the card list is person.*/
  Cards      []int /* of format 1,2,3,4 and sorted lowest -> highest */
  Rating     int8  /* 1 -> 5, 5 is best */
  NumPlayers int8
  PlayTime   int16 /* how many minutes did the game take? */
}

func writeSetRatingToDB(r *http.Request, rating *SetRating) error {
  c := appengine.NewContext(r)

  key := datastore.NewIncompleteKey(c, RatingsKindName, nil)

  _, err := datastore.Put(c, key, rating)

  return err
}

func validateAndFixupSetRating(rating *SetRating) error {
  /* for now, does nothing*/
  /* In the future, probably validates that the Cards list is valid. If we
     switch to a map written straight to the db, this will probably enforce
     a size limit on the values written */
  return nil
}

func judge(w http.ResponseWriter, r *http.Request) {
  var err error

  rating := SetRating{}
  if err = readStructFromJSONRequest(w, r, &rating); err != nil {
    serveError(w, err)
    return
  }

  if err = validateAndFixupSetRating(&rating); err != nil {
    serveError(w, err)
    return
  }

  if err = writeSetRatingToDB(r, &rating); err != nil {
    serveError(w, err)
    return
  }

  fmt.Fprintln(w, "ok")
}

/*
  this is just to show what analysis could look like
*/
func analysis(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  query := datastore.NewQuery(RatingsKindName).
    Filter("Rating >", 3)

  fmt.Fprintln(w, "These seemed to do okay:")
  for itr := query.Run(c); ; {
    rating := SetRating{}
    key, err := itr.Next(&rating)
    if err == datastore.Done {
      break
    }

    if err != nil {
      serveError(w, err)
      return
    }

    fmt.Fprintf(w, "Key=%v\nRating=%#v\n\n", key, rating)
  }
}

func matchingCard(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)

  query := datastore.NewQuery(RatingsKindName).
    Filter("Cards =", 20)

  fmt.Fprintln(w, "These had card 20:")
  for itr := query.Run(c); ; {
    rating := SetRating{}
    key, err := itr.Next(&rating)
    if err == datastore.Done {
      break
    }

    if err != nil {
      serveError(w, err)
      return
    }

    fmt.Fprintf(w, "Key=%v\nRating=%#v\n\n", key, rating)
  }
}
