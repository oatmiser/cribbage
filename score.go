package cribbage

import (
	"fmt"
	"sort"
	"strings"
)

func Score_15(h Hand) int {
	// count 15s in a five-card Hand
	/*var cardValues []int = make([]int, 0)
	for _, c := range h {
		cardValues = append(cardValues, c.ValueMax10())
	}
	// TODO subset sum to 15
	return 0
	*/
	values := make([]int, len(h))
	for i, c := range h {
		values[i] = c.ValueMax10()
	}

	points := 0
	n := len(values)
	// subset sum to 15
	// 0x00000 to 0x11111 represents the inclusion of FIVE Cards in the subset
	// iterate through the 32 cases to see if the corresponding sum is 15
	// minor optimization to start at 0x00011 (minimum 2 cards to get 15)
	for mask := 3; mask < (1 << n); mask++ {
		sum := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sum += values[i]
			}
		}
		if sum == 15 {
			points += 2
		}
	}
	return points
}

func Score_multiple(h Hand) int {
	// count pair/triple of Ranks in a five-card Hand
	var points int = 0

	count := make(map[Rank]int)
	for _, c := range h {
		count[c.Rank]++
	}
	// for rank, freq
	for _, n := range count {
		if n > 1 {
			pairs := n * (n - 1) / 2
			points += pairs * 2
		}
	}

	return points
}

func Score_run(h Hand) int {
	// hard to reason about duplicates using slices.Contains
	// frequency of each Card Rank
	/*count := make(map[Rank]int)
	for _, c := range h {
		count[c.Rank]++
	}

	// extract unique and sort
	ranksInHand := make([]Rank, 0)
	for r, count := range count {
		if count > 0 {
			ranksInHand = append(ranksInHand, r)
		}
	}
	slices.Sort(ranksInHand)

	// TODO consecutive sequences
	*/

	freq := make(map[Rank]int)
	for _, c := range h {
		freq[c.Rank]++
	}

	// extract and sort unique ranks
	ranks := make([]int, 0, len(freq))
	for r := range freq {
		ranks = append(ranks, int(r))
	}
	sort.Ints(ranks)

	maxPoints := 0
	for i := 0; i < len(ranks); i++ {
		runLen := 1
		multiplier := freq[Rank(ranks[i])]

		for j := i + 1; j < len(ranks); j++ {
			if ranks[j] == 1+ranks[j-1] {
				runLen++
				multiplier *= freq[Rank(ranks[j])]
			} else {
				break
			}
		}

		if runLen >= 3 {
			points := runLen * multiplier
			if points > maxPoints {
				maxPoints = points
			}
		}
	}

	return maxPoints
}

func Score_flush(hand Hand, cut Card, myCrib bool) int {
	/*flush := 1
	var suit Suit = h[0].Suit
	h = h[1:]
	for _, c := range h {
		if c.Suit != suit {
			break
		}
		flush++
	}

	if flush == 4 {
		if myCrib && cut.Suit == suit {
			flush++
		}
		return flush
	} else {
		return 0
	}
	*/
	if len(hand) == 0 {
		return 0
	}

	suit := hand[0].Suit
	for _, c := range hand {
		if c.Suit != suit {
			return 0
		}
	}
	// all 4 hand cards match
	if cut.Suit == suit {
		// crib or non-crib
		return 5
	}
	if myCrib {
		// only flush of 5 in the crib
		return 0
	}
	// flush of 4 or 5 for the hand
	return 4
}

func Score_nobs(h Hand, cut Card) int {
	for _, c := range h {
		if c.Rank == Jack && c.Suit == cut.Suit {
			return 1
		}
	}
	return 0
}

// ----------------------------------------------- //

type ScoreBreakdown struct {
	Fifteens int
	Pairs    int
	Runs     int
	Flush    int
	Nobs     int
	Total    int
}

func (h Hand) Score(cut Card, isCrib bool) int {
	/*var total int = 0
	all := append(h, cut)

	total += Score_15(all)
	total += Score_multiple(all)
	total += Score_run(all)
	total += Score_flush(h, cut, myCrib)
	total += Score_nobs(h, cut)

	return total
	*/
	return h.ScoreBreakdown(cut, isCrib).Total
}

// TODO make float64 for machine learn EV
func (discard Hand) HeuristicScore() int {
	// minimum possible points from this hand (i.e. 2 cards sent to Crib)
	// do not model the cut card
	// TODO machine learning: other 46 or 50 possible cut cards can influence weights?

	points := 0
	Card1, Card2 := discard[0], discard[1]

	// Only guaranteed points are by 15, pair
	// Machine learning for card weight on 15-good, multiple pair, run, flush, nobs

	// 15
	if Card1.ValueMax10()+Card2.ValueMax10() == 15 {
		points += 2
	}
	// TODO machine learning: rank card values for usefulness in getting 15

	// pair
	if Card1.Rank == Card2.Rank {
		points += 2
	}
	// TODO machine learning: minor weight on possible 3 or 4 in a row

	// TODO machine learning for weights: possible run in the crib
	diff := Card1.Value() - Card2.Value()
	if diff < 0 {
		diff = -diff
	}
	switch diff {
	case 2:
		// run of 3
		points += 0
	case 3:
		// run of 4
		points += 0
	case 4:
		// run of 5
		points += 0
	}

	// possible flush (requires 5 in crib, very rare)
	if Card1.Suit == Card2.Suit {
		points += 0
	}

	// nobs hardcode probability?
	// if you discard a Jack and other card
	if Card1.Rank == Jack || Card2.Rank == Jack {
		if Card1.Suit == Card2.Suit {
			// reduce probability???
		}
		// expected 1 point * 1/50
		// or we consider the other 4 in the keep hand.
		// Move this to DiscardOption.ExpectedValue()?
	}

	return points
}

func (h Hand) ScoreBreakdown(cut Card, isDealer bool) ScoreBreakdown {
	all := append(h, cut)
	sb := ScoreBreakdown{}
	sb.Fifteens = Score_15(all)
	sb.Pairs = Score_multiple(all)
	sb.Runs = Score_run(all)
	sb.Flush = Score_flush(h, cut, isDealer)
	sb.Nobs = Score_nobs(h, cut)
	sb.Total = sb.Fifteens + sb.Pairs + sb.Runs + sb.Flush + sb.Nobs
	return sb
}

func (sb ScoreBreakdown) Print() {
	spacer := strings.Repeat(" ", 3)
	//spacer := "  --"
	fifteens, pairs, runs, flush, nobs := sb.Fifteens, sb.Pairs, sb.Runs, sb.Flush, sb.Nobs

	if fifteens > 0 {
		fmt.Print(spacer)
		if fifteens == 2 {
			fmt.Println("15 for 2")
		} else {
			fmt.Printf("%d fifteens for %d\n", fifteens/2, fifteens)
		}
	}
	if pairs > 0 {
		fmt.Print(spacer)
		if pairs == 2 {
			fmt.Println("Pair for 2")
		} else {
			fmt.Printf("%d pairs for %d\n", pairs/2, pairs)
		}
	}
	if runs > 0 {
		fmt.Print(spacer)
		switch runs {
		case 3:
			fmt.Println("Run of 3 for 3")
		case 6:
			fmt.Println("2 runs of 3 for 6")
		case 9:
			fmt.Println("3 runs of 3 for 9")
		case 12:
			fmt.Println("4 runs of 3 for 12")
		case 4:
			fmt.Println("Run of 4 for 4")
		case 8:
			fmt.Println("2 runs of 4 for 8")
		case 5:
			fmt.Println("Run of 5 for 5")
		case 10:
			fmt.Println("2 runs of 5 for 10")
		default:
			fmt.Println("TODO ScoreBreakdown")
		}
		//fmt.Printf("%d runs of X for %d\n", runs/2, runs)
	}
	if flush > 0 {
		fmt.Printf(spacer+"Flush for %d\n", flush)
	}
	if nobs > 0 {
		fmt.Println(spacer + "Nobs for 1")
	}
}
