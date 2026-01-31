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
	options := p.Hand.Split(4)
	best := OptimalDiscard(options, isDealer)
	discard = best.Discard
	keep = best.Keep

	// copy of Hand will be emptied during Pegging
	p.Hand = keep
	p.PegHand = make(Hand, len(keep))
	copy(p.PegHand, keep)
	return
}

func (p *ComputerPlayer) PlayPegCard(s PegState) (cardToPlay Card, passed bool) {
	best, ok := OptimalPegging(s, p.PegHand)
	if ok {
		cardToPlay = best
		passed = false
		return
	}

	// else no optimal detected, play the first valid one
	for i, card := range p.PegHand {
		if card.ValueMax10() <= 31-s.Sum {
			p.PegHand = slices.Delete(p.PegHand, i, i+1)
			cardToPlay = card
			passed = false
			return
		}
	}
	// no valid card and say Go/pass
	cardToPlay = Card{}
	passed = true
	return
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
	// Function in Player Interface, only needed for HumanPlayer
	//return
}

func (p *ComputerPlayer) EmptyPegHand() bool {
	return len(p.PegHand) == 0
}
