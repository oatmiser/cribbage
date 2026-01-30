package cribbage

import (
	"fmt"
	"strings"
)

type DiscardOption struct {
	Keep    Hand
	Discard Hand
	//ExpectedValue float64
}

func (opt DiscardOption) ScoreRange(isDealer bool) (int, int) {
	// best hand is 29 points
	min, max := 29, 0
	knownRemaining := difference(Hand(NewDeck()), append(opt.Keep, opt.Discard...))

	for _, cut := range knownRemaining {
		// Score during Show, crib is false
		possiblePoints := opt.Keep.Score(cut, false)
		if isDealer {
			possiblePoints += opt.Discard.HeuristicScore()
		} else {
			// effectively lost points by giving them to opponent
			possiblePoints -= opt.Discard.HeuristicScore()
		}

		if possiblePoints < min {
			min = possiblePoints
		}
		if possiblePoints > max {
			max = possiblePoints
		}
	}
	return min, max
}

func DiscardAnalysis(options []DiscardOption, isDealer bool) {
	// every possible Keep/Discard and their average score
	msg := "All Possible Discards"
	if isDealer {
		msg += " (Your Crib)"
	} else {
		msg += " (Opponent's Crib)"
	}
	fmt.Println(msg)
	fmt.Println(strings.Repeat("-", len(msg)))
	for i, opt := range options {
		fmt.Printf("Option #%d\n", i+1)
		ev := opt.ExpectedValue(isDealer)
		min, max := opt.ScoreRange(isDealer)
		fmt.Printf("Hand: %s\nCrib: %s\n", opt.Keep, opt.Discard)
		//fmt.Printf("Average points: %f\nMin score: %d\nMax score: %d\n\n", ev, min, max)
		fmt.Printf("Average points: %f\n", ev)
		fmt.Printf("Score min, max = %d, %d\n\n", min, max)
	}
	fmt.Println(strings.Repeat("-", len(msg)))
	//PrintOptimal(options, isDealer)
}

func OptimalDiscard(options []DiscardOption, isDealer bool) DiscardOption {
	var bestOption DiscardOption = options[0]
	var bestEV float64 = 0.0

	// every possible Keep/Discard and their average score
	for _, option := range options {
		ev := option.ExpectedValue(isDealer)
		if ev > bestEV {
			bestEV = ev
			bestOption = option
		}
	}
	return bestOption
}

func PrintOptimal(options []DiscardOption, isDealer bool) {
	optimal := OptimalDiscard(options, isDealer)
	var player string
	if isDealer {
		player = "dealer"
	} else {
		player = "pone"
	}

	fmt.Printf("As the %s, your optimal discard is %s\n", player, optimal.Discard)
	//fmt.Printf(" (EV %f)\n\n", optimal.ExpectedValue(isDealer))
	//fmt.Printf("Expected value is %f", optimal.ExpectedValue(isDealer))
	min, _ := optimal.ScoreRange(isDealer)
	//fmt.Printf("; you will get at least %d points.\n", min)
	// TODO needs to model cut card...
	fmt.Printf("Your remaining Hand will get at least %d points", min)
	fmt.Printf(" and the expected value is %f\n", optimal.ExpectedValue(isDealer))
}

func (opt DiscardOption) ExpectedValue(isDealer bool) float64 {
	// Hand and Deck are []Card
	var sumPoints int = 0
	var minimumCrib int = 0
	knownRemaining := difference(Hand(NewDeck()), append(opt.Keep, opt.Discard...))

	// For every possible cut card, the Player would get that score during Show
	for _, cut := range knownRemaining {
		// Score during Show, crib is false
		sumPoints += opt.Keep.Score(cut, false)
	}
	// Expected Points is the average of Show points from every possible cut card
	// Anything in our Keep / Discard is not a cut, and we do not know the opponent's hand
	count := float64(len(knownRemaining)) // 46
	expectedShow := float64(sumPoints) / count

	// TODO assuming that best Show points is the best Keep/Discard is WRONG

	// if isDealer we count some points from discard
	if isDealer {
		// Score during Show with Crib (only 2 cards known)
		// 52 choose 2 is 1326, 46... is 1035
		minimumCrib += opt.Discard.HeuristicScore()
	} else {
		minimumCrib -= opt.Discard.HeuristicScore()
	}

	// TODO model opponent
	// avoid sum (sometimes?) to 5, 10, 21

	return expectedShow + float64(minimumCrib)
}
