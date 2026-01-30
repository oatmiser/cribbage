package cribbage

import (
	"fmt"
	"math/rand/v2"
	"slices"
)

type ComputerPlayer struct {
	Name    string
	Hand    Hand
	PegHand Hand
	Points  int
}

func (p *ComputerPlayer) String() string {
	return p.Name
}

func (p *ComputerPlayer) GetName() string {
	return p.Name
}

func (p *ComputerPlayer) GetScore() int {
	return p.Points
}

func (p *ComputerPlayer) GetHand() Hand {
	return p.Hand
}

func (p *ComputerPlayer) SetHand(h Hand) {
	p.Hand = h
}

func (p *ComputerPlayer) AddPoints(n int) int {
	p.Points += n
	return p.Points
}

func (p *ComputerPlayer) Discard(isDealer bool) (discard Hand, keep Hand) {
	//h := p.Hand
	options := p.Hand.Split(4)
	bestOption := OptimalDiscard(options, isDealer)
	//discard = h[:2]
	//keep = h[2:]
	discard = bestOption.Discard
	keep = bestOption.Keep

	p.Hand = keep
	// copy of Hand will be emptied during Pegging
	p.PegHand = make(Hand, len(keep))
	copy(p.PegHand, keep)
	return
}

func (p *ComputerPlayer) PeggingHeuristic(s PegState) Card {
	// For every playable card in Hand...
	// Immediate gain
	// Target 15, 31, pair, run

	// Subtract risk score (enables opponent points)
	// Avoid 5, 10, 21, opponent run

	// Predictive reasoning
	// Allow a card if optimal opponent behavior can be used
	// e.g. play 7, expect 8, play 9 (+3 run > +2 15 for opponent)

	// Position from 15 or 31
	// Force opponent GO

	return Card{}
}

func (p *ComputerPlayer) PlayPegCard(s PegState) (c Card, passed bool) {
	// send first valid card in Player's hand
	// TODO: consider run, in a row, 15, 31, ?predict other player?
	best, ok := OptimalPegging(s, p.PegHand)
	if ok {
		c = best
		passed = false
		return
	}

	// else no optimal detected, play the first valid one
	for i, card := range p.PegHand {
		if card.ValueMax10() <= 31-s.Sum {
			p.PegHand = slices.Delete(p.PegHand, i, i+1)
			return card, false
		}
	}
	// no valid card and say Go/pass
	return Card{}, true
}

func (p *ComputerPlayer) DrawCard() int {
	return rand.IntN(52)
}

func (p *ComputerPlayer) CountHand(cut Card, isCrib bool) int {
	points := p.Hand.ScoreBreakdown(cut, isCrib)
	if isCrib {
		fmt.Printf("%s (Crib): %s", p.Name, p.Hand)
	} else {
		fmt.Printf("%s: %s", p.Name, p.Hand)
	}
	fmt.Printf(" (%d points)\n", points.Total)
	points.Print()
	return points.Total
}

func (p *ComputerPlayer) EnterToContinue() {
	// Function in Interface needed for HumanPlayer
	//return
}

func (p *ComputerPlayer) EmptyPegHand() bool {
	return len(p.PegHand) == 0
}
