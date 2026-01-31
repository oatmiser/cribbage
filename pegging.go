package cribbage

// File contains loop for Pegging round and functions for scoring Pegging points

import (
	"fmt"
	"strings"
)

type PegState struct {
	Sum        int     // pegging up to 31
	Turn       int     // index of Player
	LastPlayer int     // last player to place a card
	CardPile   Hand    // current Player cards (add to <=31)
	Passed     [2]bool // player said "Go"
	PileNum    int     // 1-3 piles of <=31 per Pegging round
}

// Players try to place all of their cards on the pile
// Once the score goes above 31, start a new pile
func (game *Game) StartPegging() {
	dealer := game.Dealer
	players := game.Players
	linebreak := strings.Repeat("-", 30)

	state := PegState{
		Sum:      0,
		Turn:     1 - dealer, // pone places first card
		CardPile: make([]Card, 0),
	}

	fmt.Printf("(Pegging Pile 1)\n")
	for !EmptyHands(players) {
		//fmt.Printf("(Pegging Pile %d)\n", state.PileNum+1)
		// assume another skip if Player previously passed
		if state.Passed[state.Turn] {
			//fmt.Printf("%s says GO\n\n", players[state.Turn])
			state.Turn = 1 - state.Turn
			continue
		}

		fmt.Printf("Sum: %d\n", state.Sum)
		for _, c := range state.CardPile {
			fmt.Printf("%s ", c)
		}
		fmt.Println("[?]")

		card, passed := players[state.Turn].PlayPegCard(state)
		if passed {
			fmt.Printf("%s says GO", players[state.Turn])
			state.Passed[state.Turn] = true
		} else {
			fmt.Printf("%s plays %s", players[state.Turn], card)
			points, comment := ScorePeggingPlay(state, card)
			state.AddCard(card)
			if points > 0 {
				fmt.Print(comment)
				game.AddPoints(state.Turn, points)
				if game.GameWon {
					game.CelebrateWinner(state.Turn)
					return
				}
			}
		}

		if state.ShouldReset() {
			// give points for last card (31 is included in ScorePeggingPlay)
			if state.Sum != 31 {
				fmt.Println()
				/*fmt.Printf("\nSum: %d\n", state.Sum)
				for _, c := range state.CardPile {
					fmt.Printf("%s ", c)
				}
				*/
				fmt.Printf("\n%s scores +1 [Last Card]", players[state.LastPlayer])
				game.AddPoints(state.LastPlayer, 1)
				if game.GameWon {
					game.CelebrateWinner(state.LastPlayer)
					return
				}
			}
			fmt.Printf("\n\n")
			state.Reset()
			if !EmptyHands(players) {
				fmt.Printf("(Pegging Pile %d)\n", state.PileNum+1)
			}
		} else {
			fmt.Printf("\n%s\n", linebreak)
		}

		state.Turn = 1 - state.Turn
	}

	//fmt.Printf("\nSum: %d\n", state.Sum)
	/*for _, c := range state.CardPile {
		fmt.Printf("%s ", c)
	}
	*/
	fmt.Println(state.CardPile)
	fmt.Printf("\nAll cards have been played!\n")

	if state.Sum != 0 {
		// the loop ended after both hands are empty
		// sum can be zero after 31, or if both Players Go
		fmt.Printf("%s scores +1 [Last Card]\n\n", players[state.LastPlayer])
		game.AddPoints(state.LastPlayer, 1)
		if game.GameWon {
			game.CelebrateWinner(state.LastPlayer)
			return
		}
	}
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
				// continue in outer loop, where the next (smaller)
				// topN might make a valid run without this duplicate
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
	// state was not updated with Card in the caller yet
	s.AddCard(c)

	points = 0

	if s.Sum == 15 {
		msg = "  +2  [15]"
		points += 2
	}
	if s.Sum == 31 {
		msg = "  +2  [31]"
		points += 2
	}

	pairPoints := ScorePegPairs(s.CardPile)
	if pairPoints == 2 {
		msg += "  +2  [Pair]"
	} else if pairPoints > 1 {
		// TODO say 3 in a row etc
		msg += fmt.Sprintf("  +%d  [%d pairs]", pairPoints, pairPoints/2)
	}
	points += pairPoints

	runPoints := ScorePegRuns(s.CardPile)
	if runPoints > 0 {
		msg += fmt.Sprintf("  +%d  [Run of %d]", runPoints, runPoints)
	}
	points += runPoints

	// 1 point can also be given for "last card of a pile" inside StartPegging
	return
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
	return players[0].EmptyPegHand() && players[1].EmptyPegHand()
}

func OptimalPegging(state PegState, hand Hand) (Card, bool) {
	// TODO opponent modeling and more, examples:
	// avoid allowing run of 3 unless you can make run of 4
	// do not do 7-8 (+2 15 allows +3 run of 3)
	// avoid 5, 10, 21
	max := 0
	var best Card
	found := false

	for _, card := range hand {
		// cannot ever play a card that is too big
		if card.ValueMax10()+state.Sum > 31 {
			continue
		}
		points, _ := ScorePeggingPlay(state, card)
		if points > max {
			max = points
			best = card
			found = true
		}
	}
	return best, found
}
