package cribbage

import (
	"fmt"
	"strconv"
)

// Enumerations (Typed Constants)
type Rank int

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

type Color int

const (
	Black Color = iota
	Red
)

// Card Data Structure
type Card struct {
	Rank  Rank
	Suit  Suit
	Color Color
}

func (c Card) ValueMax10() int {
	switch c.Rank {
	case Jack, Queen, King:
		return 10
	default:
		return int(c.Rank)
	}
}

// Value used for runs
func (c Card) Value() int {
	return int(c.Rank)
}

func (c Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

// String Conversions for CLI
func (r Rank) String() string {
	switch r {
	case Ace:
		return "A"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		//return string(rune('0') + rune(r))
		return strconv.Itoa(int(r))
	}
}

func (s Suit) String() string {
	switch s {
	case Clubs:
		return "♣"
	case Diamonds:
		return "♦"
	case Hearts:
		return "♥"
	case Spades:
		return "♠"
	default:
		return "?"
	}
}
