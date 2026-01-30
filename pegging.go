package cribbage

import (
	"fmt"
)

type PegState struct {
	Sum        int // pegging up to 31
	Turn       int
	LastPlayer int // last player to place a card
	CardPile   Hand
	Passed     [2]bool // player said "Go"
	PileNum    int     // 1-3 piles of <=31 per Pegging round
}

func TrailingMultiple(pegStack Hand) int {
	if len(pegStack) < 2 {
		return 0
	}
	last := pegStack[len(pegStack)-1].Rank
	count := 1

	for i := len(pegStack) - 2; i >= 0; i-- {
		if pegStack[i].Rank != last {
			break
		}
		count++
	}

	return count
}

func ScorePegPairs(stack Hand) int {
	switch TrailingMultiple(stack) {
	case 2:
		return 2
	case 3:
		return 6
	case 4:
		return 12
	default:
		return 0
	}
}

func ScorePegRuns(stack Hand) int {
	if len(stack) < 3 {
		return 0
	}

	for n := len(stack); n >= 3; n-- {
		// first run to be witnessed will have the most points
		topN := stack[len(stack)-n:]
		// break if a duplicate Rank is seen
		seen := make(map[int]bool)
		// run from X to Y will be of length Y-X+1...
		// run can be counted only in the topN cards (max-min+1==n)
		max := Card{Rank: Ace}.Value()
		min := Card{Rank: King}.Value()

		for _, card := range topN {
			r := int(card.Rank)
			if seen[r] {
				// continue in outer loop...
				// next topN (smaller) may have a valid run without the duplicate
				goto next
			}
			seen[r] = true

			if r < min {
				min = r
			}
			if r > max {
				max = r
			}
		}

		if max-min+1 == n {
			return n
		}

	next:
	}
	return 0
}

// amount of points from placing such a Card on the pegging pile
// with the string to append in the terminal output
func ScorePeggingPlay(s PegState, c Card) (points int, msg string) {
	// TODO bad design!!!
	// state must be ALREADY changed with card c
	//fixed
	s.AddCard(c)

	/*spacer := strings.Repeat(" ", 4)
	spacer = "  --"
	spacer = "\n --"
	*/
	points = 0

	if s.Sum == 15 {
		//fmt.Printf("%s15 for 2\n", spacer)
		msg = "  +2  [15]"
		points += 2
	}
	if s.Sum == 31 {
		//fmt.Printf("%s31 for 2\n", spacer)
		msg = "  +2  [31]"
		points += 2
	}

	pairPoints := ScorePegPairs(s.CardPile)
	if pairPoints == 2 {
		//fmt.Printf("%sPair for 2\n", spacer)
		msg += "  +2  [Pair]"
	} else if pairPoints > 1 {
		//fmt.Printf("%s%d pairs for %d\n", spacer, pairPoints/2, pairPoints)
		// todo say 3 in a row etc?
		msg += fmt.Sprintf("  +%d  [%d pairs]", pairPoints, pairPoints/2)
	}
	points += pairPoints

	runPoints := ScorePegRuns(s.CardPile)
	if runPoints > 0 {
		//fmt.Printf("%sRun of %d for %d\n", spacer, runPoints, runPoints)
		msg += fmt.Sprintf("  +%d  [Run of %d]", runPoints, runPoints)
	}
	points += runPoints

	// Separate point is given for last card of a pile in StartPegging
	return //points, msg
}

func (s *PegState) AddCard(c Card) {
	if s.Sum+c.ValueMax10() > 31 {
		fmt.Printf("TODO over 31 in AddCard")
	}
	s.Sum += c.ValueMax10()
	s.CardPile = append(s.CardPile, c)
	s.LastPlayer = s.Turn
}

func (s *PegState) ShouldReset() bool {
	if s.Sum == 31 {
		return true
	}
	return s.Passed[0] && s.Passed[1]
}

func (s *PegState) Reset() {
	s.Sum = 0
	s.Passed = [2]bool{}
	s.Turn = s.LastPlayer
	s.CardPile = make([]Card, 0)
	s.PileNum++
}

func EmptyHands(players [2]Player) bool {
	/*h1 := players[0].GetHand()
	h2 := players[1].GetHand()
	return len(h1) == 0 && len(h2) == 0
	*/
	return players[0].EmptyPegHand() && players[1].EmptyPegHand()
}

func OptimalPegging(state PegState, hand Hand) (Card, bool) {
	// TODO opponent modeling and more
	// avoid allowing run of 3 unless you can make run of 4
	// avoid taking run of 4 unless you can make run of 4/5
	// e.g. do not do 7-8 (+2 15 but allows +3 run of 3)
	// avoid 5, 10, 21
	max := 0
	var best Card
	found := false

	for _, card := range hand {
		if card.ValueMax10()+state.Sum > 31 {
			// cannot ever play a card that is too big
			//fmt.Printf("Cannot play %s\n", card)
			continue
		}

		points, _ := ScorePeggingPlay(state, card)
		//fmt.Printf("%s gets %d points\n", card, points)
		if points > max {
			max = points
			best = card
			found = true
		}
	}
	return best, found
}
