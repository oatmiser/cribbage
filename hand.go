package cribbage

import (
	"fmt"
	"strings"
)

// Player Hand from Deck
type Hand []Card

func (h Hand) Print(name string) {
	fmt.Printf("%s: ", name)
	for i := range len(h) {
		//fmt.Print(h[i].Rank.String())
		//fmt.Print(h[i].Suit.String() + "  ")
		fmt.Printf("%s ", h[i])
	}
	fmt.Println()
}

func (h Hand) String() string {
	var ret strings.Builder
	for _, card := range h {
		//ret.WriteString(fmt.Sprintf("%s ", card))
		fmt.Fprintf(&ret, "%s ", card)
	}
	// remove last space
	//return ret.String()[:2*len(h)-1]
	return strings.TrimRight(ret.String(), " ")
}

func (h Hand) Choose(k int) []Hand {
	var result []Hand
	var helper func(start int, current Hand)

	helper = func(start int, current Hand) {
		if len(current) == k {
			// use a copy to avoid aliasing
			hh := make(Hand, k)
			copy(hh, current)
			result = append(result, hh)
			return
		}
		// 0..2, {}
		// 1..1, {h[0]}
		for i := start; i <= len(h)-(k-len(current)); i++ {
			// 2 - len(current)
			helper(i+1, append(current, h[i]))
		}
	}

	helper(0, Hand{})
	return result
}

// generate all ways to keep 4 and discard 2 at each round
func (hand Hand) Split(out int) []DiscardOption {
	in := len(hand)
	if in != 6 {
		fmt.Printf("Not a dealt hand of 6 cards!")
		return nil
	}
	if out != 4 {
		fmt.Printf("Cribbage hand must be 4 cards!")
		return nil
	}

	var options []DiscardOption
	keepList := hand.Choose(out)
	// check each subslice (keep pile) from 6 choose 4 on the Hand
	for _, keep := range keepList {
		// figure out the 2 discarded cards
		discard := difference(hand, keep)
		// Expected Value of 0 is unused
		// It is ignored or separately set by the caller of Split
		options = append(options, DiscardOption{keep, discard})
	}

	return options
}

// Set Difference between a Hand6 and Hand4
func difference(full, subset Hand) Hand {
	marked := make(map[Card]int)
	// mark all cards in the subset (player's keep pile)
	for _, card := range subset {
		marked[card]++
	}

	// compare to all cards in the original (player's dealt hand)
	var diff Hand
	for _, card := range full {
		if marked[card] > 0 {
			marked[card]--
		} else {
			diff = append(diff, card)
		}
	}
	return diff
}
