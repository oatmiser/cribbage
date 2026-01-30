package cribbage

/*
Computer player will learn Cribbage gameplay and the importance of
e.g. pegging pile sum, last card for pair/run, and controlling piles
using the distance from 31 and hypothesized opponent hand
*/
type PegStateFeatures struct {
	sum            int      // pile is at 0 up to 31
	state          PegState // too complicated for feature?
	lastCardRank   Rank
	runLength      int
	cardsRemaining Hand
	opponentCards  Hand // estimated worst hand
	opponentPassed bool
}

func RiskScore(sum int) int {
	switch sum {
	case 5, 10, 21:
		return -3 // opponent can likely get 15 or 31
	case 14, 20, 26:
		return -1 // opponent can get 15, 31, last card
	default:
		return 0
	}
}

func PredictivePenalty(sum int) int {
	// How likely opponent can score?
	if 22 <= sum && sum <= 26 {
		return -2 // many ways to hit 31
	}
	if 10 <= sum && sum < 15 {
		return -1 // many ways to hit 15
	}
	return 0
}

func PositionScore(sum int, handSize int) int {
	if sum == 31 {
		return 3
	}
	if sum >= 27 {
		return 2
	}
	if handSize == 1 {
		return 2 // last-card pressure
	}
	return 0
}
