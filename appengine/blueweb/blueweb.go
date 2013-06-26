package blueweb

import (
  "appengine"
  "appengine/datastore"
  "crypto/rand"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "math/big"
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

type Get10Options struct {
  NoAttack               bool
  NoAttackWithoutDefense bool
  SetsAvailable          []int
  ForbiddenCards         []int
}

func readStructFromJSONRequest(w http.ResponseWriter, r *http.Request, readInto interface{}) error {
  // TODO: we should probably check that this is a POST and that it's a request of type text/json
  jsonRaw, _ := ioutil.ReadAll(r.Body)
  return json.Unmarshal(jsonRaw, readInto)
}

func generateRandomNumbers(maxVal int, c chan int64) {
  max := big.NewInt(int64(maxVal))

  for {
    bigInt,_ := rand.Int(rand.Reader, max)
    c <- bigInt.Int64()
  }

} 

func generateRandomCards(options *Get10Options, c chan int) {
/*
  The general strategy here is to generate a list with all the cards in it,
  then randomly pick indexes out of - we then mark that item in the list
  as used. If we happen to pick an already used card, we will just keep
  picking until we pick an unused card. Since the number of cards needed
  should be significantly smaller than the size of the deck, collisions shouldn't
  be an issue.

  I think that will end up being something like O(cards_needed), so
  somewhat close to 10.

  We could have generated the deck, then sorted it, but we think that the
  deck size will always be at least 2x number of cards needed. sorting the deck is guaranteed 
  to be O(nlogn), which I expect to be much larger than O(cards_needed), but the pathological
  case is an infinite loop. We think the pathological case won't happen that often,
  this is good enough for now.
*/


  // todo: use the actual set of cards, and pull things out based on the options
  cards := []int{1, 2, 3, 7, 8, 6, 5, 20, 18, 19, 14, 12, 13}

  randomNumbers := make(chan int64)
  go generateRandomNumbers(len(cards), randomNumbers)
  for numGenerated := 0; numGenerated < len(cards); numGenerated += 1 {
    curCard := <- randomNumbers
    for cards[curCard] == 0 {  
      curCard = <- randomNumbers
    }
    c <- cards[curCard]
    cards[curCard] = 0
    // todo: make it so we don't die in the pathological case
    // one idea is that when numGenerated > len(cards)/2, 
    // clear all the used values out of cards
  }

  close (c)
}

// this is where all the magic happens for generating the 10 cards
func generateCards(options *Get10Options) []int {

  c := make(chan int)
  go generateRandomCards(options, c)
  
  cards := make([]int, 10)

  for i := 0; i < 10; i++ {
    cards[i] = <- c
  }

  return cards
}

func generateGet10Response(cards []int, w http.ResponseWriter) {
  b, _ := json.Marshal(cards)
  fmt.Fprint(w, string(b))
}

func get10(w http.ResponseWriter, r *http.Request) {
  var err error

  options := Get10Options{}
  if err = readStructFromJSONRequest(w, r, &options); err != nil {
    serveError(w, err)
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
