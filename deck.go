package cribbage

import (
	"math/rand"
	"time"
)

// Deck of Cards
type Deck []Card

func NewDeck() Deck {
	// len 0 and cap 52
	deck := make(Deck, 0, 52)
	// Clubs 0, Diamonds 1, Hearts 2, Spades 3
	for suit := Clubs; suit <= Spades; suit++ {
		color := Black
		if suit == Diamonds || suit == Hearts {
			color = Red
		}
		for rank := Ace; rank <= King; rank++ {
			deck = append(deck, Card{rank, suit, color})
		}
	}

	return deck
}

func (d Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	// TODO implement my own Fisher-Yates shuffle
	rand.Shuffle(len(d), func(i, j int) {
		d[i], d[j] = d[j], d[i]
	})
}

func Deal(deck Deck, count int) (Hand, Hand, Deck) {
	//count := 6
	p1 := make(Hand, 0, count)
	p2 := make(Hand, 0, count)
	for i := 0; i < count; i++ {
		//for i := range count {
		p1 = append(p1, deck[0])
		p2 = append(p2, deck[1])
		deck = deck[2:]
	}
	return p1, p2, deck
}

func ShowCardDeal(p1, p2 Player) {
	p1.GetHand().Print(p1.GetName())
	p2.GetHand().Print(p2.GetName())
}
