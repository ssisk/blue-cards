package generate

import (
  "blueshared"
  "crypto/rand"
  "math/big"
)

func generateRandomNumbers(maxVal int, c chan int64) {
  max := big.NewInt(int64(maxVal))

  for {
    bigInt,_ := rand.Int(rand.Reader, max)
    c <- bigInt.Int64() 
  }

} 
  
func generateRandomCards(options *blueshared.Get10Options, c chan int) {
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
func CardSet(options *blueshared.Get10Options) []int {

  c := make(chan int)
  go generateRandomCards(options, c)
  
  cards := make([]int, 10)

  for i := 0; i < 10; i++ {
    cards[i] = <- c
  }

  return cards
}