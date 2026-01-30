package cribbage

import (
	"testing"
)

func makeStack(ranks ...Rank) Hand {
	h := Hand{}
	for _, r := range ranks {
		h = append(h, Card{Rank: r})
	}
	return h
}

func TestPegging_15_31(t *testing.T) {
	tests := []struct {
		name  string
		sum   int
		stack Hand
		want  int
	}{
		{"Peg 15", 15, makeStack(Five, Ten), 2},
		{"Peg 31", 31, makeStack(Ten, Ten, Nine, Two), 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PegState{
				Sum:      tt.sum,
				CardPile: tt.stack,
			}
			if got, _ := ScorePeggingPlay(s, tt.stack[len(tt.stack)-1]); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPegging_Pairs(t *testing.T) {
	tests := []struct {
		name  string
		stack Hand
		want  int
	}{
		{"TwoKind", makeStack(Six, Six), 2},
		{"ThreeKind", makeStack(Six, Six, Six), 6},
		{"FourKind", makeStack(Six, Six, Six, Six), 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PegState{CardPile: tt.stack}
			//if got := ScorePeggingPlay(&s, tt.stack[len(tt.stack)-1]); got != tt.want {
			if got := ScorePegPairs(s.CardPile); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestPegging_Runs(t *testing.T) {
	tests := []struct {
		name  string
		stack Hand
		want  int
	}{
		{"Run3", makeStack(Five, Three, Four), 3},
		{"Run4", makeStack(Five, Four, Three, Two), 4},
		{"Unordered5", makeStack(Ace, Four, Five, Two, Three), 5},
		{"NoRun", makeStack(Five, Four, Two), 0},
		{"BrokenRun", makeStack(Five, Four, Three, Three), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := PegState{CardPile: tt.stack}
			//if got := ScorePeggingPlay(&s, tt.stack[len(tt.stack)-1]); got != tt.want {
			if got := ScorePegRuns(s.CardPile); got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

// happens outside of ScorePeggingPlay()
func TestPegging_GoGo(t *testing.T) {
	p := &ComputerPlayer{Name: "P1", Points: 0}
	state := PegState{
		Sum:        29,
		LastPlayer: 0,
		Passed:     [2]bool{true, true},
	}

	//p.AddPoints(0)
	if state.ShouldReset() {
		if state.Sum != 31 {
			p.AddPoints(1)
		}
	}

	if p.Points != 1 {
		t.Fatalf("got %d, want 1", p.Points)
	}
}
