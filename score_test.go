package cribbage

import (
	"testing"
)

func TestCountFifteens(t *testing.T) {
	tests := []struct {
		name string
		hand Hand
		want int
	}{
		{
			name: "single fifteen",
			hand: Hand{
				{Rank: Five}, {Rank: Ten}, {Rank: Nine}, {Rank: Nine}, {Rank: Nine},
			},
			want: 2,
		},
		{
			name: "multiple fifteens",
			hand: Hand{
				{Rank: Five}, {Rank: Five}, {Rank: Five}, {Rank: Jack}, {Rank: Queen},
			},
			want: 14,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Score_15(tt.hand); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCountPairs(t *testing.T) {
	tests := []struct {
		name string
		hand Hand
		want int
	}{
		{
			name: "one pair",
			hand: Hand{
				{Rank: Ace}, {Rank: Ace}, {Rank: Two}, {Rank: Three}, {Rank: Four},
			},
			want: 2,
		},
		{
			name: "two pairs",
			hand: Hand{
				{Rank: Ace}, {Rank: Ace}, {Rank: Two}, {Rank: Queen}, {Rank: Queen},
			},
			want: 4,
		},
		{
			name: "three of a kind",
			hand: Hand{
				{Rank: Five}, {Rank: Five}, {Rank: Five}, {Rank: Six}, {Rank: Seven},
			},
			want: 6,
		},
		{
			name: "four of a kind",
			hand: Hand{
				{Rank: Seven}, {Rank: Seven}, {Rank: Seven}, {Rank: Seven}, {Rank: Ace},
			},
			want: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Score_multiple(tt.hand); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCountRuns(t *testing.T) {
	tests := []struct {
		name string
		hand Hand
		want int
	}{
		{
			name: "run of 3",
			hand: Hand{
				{Rank: Three}, {Rank: Four}, {Rank: Five}, {Rank: King}, {Rank: Ace},
			},
			want: 3,
		},
		{
			name: "double run of 4",
			hand: Hand{
				{Rank: Three}, {Rank: Three}, {Rank: Four}, {Rank: Five}, {Rank: Six},
			},
			want: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Score_run(tt.hand); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCountFlush(t *testing.T) {
	hand := Hand{
		{Suit: Hearts}, {Suit: Hearts}, {Suit: Hearts}, {Suit: Hearts},
	}
	cut := Card{Suit: Hearts}

	if got := Score_flush(hand, cut, true); got != 5 {
		t.Fatalf("got %d, want 5", got)
	}
}

func TestCountNobs(t *testing.T) {
	hand := Hand{
		{Rank: Jack, Suit: Spades},
	}
	cut := Card{Suit: Spades}

	if got := Score_nobs(hand, cut); got != 1 {
		t.Fatalf("got %d, want 1", got)
	}
}
