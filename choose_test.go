package cribbage

import (
	"testing"
)

func testHand(n int) Hand {
	h := make(Hand, 0, n)
	for i := range n {
		h = append(h, Card{Rank(i + 1), Clubs, Black})
	}
	return h
}

// Test that the recursive helper in Choose returns with
// the correct number of combinations.
func TestChoose_Count(t *testing.T) {
	tests := []struct {
		name     string
		handSize int
		k        int
		want     int
	}{
		{"6 choose 4", 6, 4, 15},
		{"5 choose 3", 5, 3, 10},
		{"4 choose 2", 4, 2, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := testHand(tt.handSize)
			got := len(h.Choose(tt.k))
			if got != tt.want {
				t.Fatalf("Choose(%d) = %d, want %d", tt.k, got, tt.want)
			}
		})
	}
}

// Test that Choose returns with 4 cards in every possible combination.
func TestChoose_CombinationSize(t *testing.T) {
	h := testHand(6)
	combinations := h.Choose(4)

	for i, cards := range combinations {
		if len(cards) != 4 {
			t.Fatalf("combination %d has size %d, want 4", i, len(cards))
		}
	}
}

// Test that Choose returns with no Card duplicates within any combination.
func TestChoose_NoDuplicates(t *testing.T) {
	h := testHand(6)
	combinations := h.Choose(4)

	for i, cards := range combinations {
		seen := make(map[Card]bool)
		for _, card := range cards {
			if seen[card] {
				t.Fatalf("duplicate card in combination %d: %+v", i, card)
			}
			seen[card] = true
		}
	}
}

// Test that Choose returns only with combinations that are subsets
// of the dealt (6 Card) hand.
func TestChoose_SubsetOfOriginal(t *testing.T) {
	h := testHand(6)
	combinations := h.Choose(4)

	original := make(map[Card]bool)
	// mark each Card in the Hand
	for _, c := range h {
		original[c] = true
	}
	for i, cards := range combinations {
		for _, card := range cards {
			if !original[card] {
				t.Fatalf("combination %d contains card not in hand: %+v", i, card)
			}
		}
	}
}
