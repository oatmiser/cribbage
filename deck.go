package cribbage

import (
	"math/rand"
)

// Deck of Cards
type Deck []Card

func NewDeck() Deck {
	// len 0 and card capacity 52
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

func (d Deck) Shuffle(rng *rand.Rand) {
	// Fisher-Yates shuffle by swapping elements
	n := len(d)

	for i := n - 1; i > 0; i-- {
		j := rng.Intn(i + 1) // j from [0,i]
		d[i], d[j] = d[j], d[i]
	}
}

func Deal(deck Deck, count int) (Hand, Hand, Deck) {
	p1 := make(Hand, 0, count)
	p2 := make(Hand, 0, count)
	for i := 0; i < count; i++ {
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
