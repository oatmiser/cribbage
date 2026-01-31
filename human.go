package cribbage

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"slices"
	"strconv"
	"strings"
)

type HumanPlayer struct {
	Name    string
	Hand    Hand
	PegHand Hand
	Points  int
}

func (p *HumanPlayer) String() string {
	return p.Name
}

func (p *HumanPlayer) GetName() string {
	return p.Name
}

func (p *HumanPlayer) GetHand() Hand {
	return p.Hand
}

func (p *HumanPlayer) SetHand(h Hand) {
	p.Hand = h
}

func (p *HumanPlayer) AddPoints(n int) int {
	p.Points += n
	return p.Points
}

func (p *HumanPlayer) GetScore() int {
	return p.Points
}

func (p *HumanPlayer) Discard(isDealer bool) (discard Hand, keep Hand) {
	// Player's Hand was created by SetHand and PegHand is CURRENTLY NIL
	// Copy the Keep (4 cards) to PegHand so it can be emptied during Pegging
	dealtHand := p.Hand // 6 cards

	if isDealer {
		fmt.Println("Select 2 cards to send to your Crib.")
	} else {
		fmt.Println("Select 2 cards to send to the opponent's Crib.")
	}
	fmt.Println("Say 'h' for a Hint on optimal selection")
	fmt.Println("Examples: '1 2', '4 1', '6 2'")
	PromptIndices(dealtHand)

	//var input string
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Select Cards with 2 indices separate by a space: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch strings.ToLower(input) {
		case "hint", "help", "h":
			PrintOptimal(p.Hand.Split(4), isDealer)
			PromptIndices(dealtHand)
			continue
		case "all", "analyze", "a":
			DiscardAnalysis(p.Hand.Split(4), isDealer)
			PromptIndices(dealtHand)
			continue
		case "optimal", "opt", "o", "win", "w":
			PrintOptimal(p.Hand.Split(4), isDealer)
			fmt.Println()
			p.EnterToContinue()

			optimal := OptimalDiscard(p.Hand.Split(4), isDealer)
			p.Hand = optimal.Keep
			p.PegHand = make(Hand, len(p.Hand))
			copy(p.PegHand, p.Hand)
			return optimal.Discard, optimal.Keep
		}

		// Try to parse all other cases of user input as "<int> <int>"
		fields := strings.Fields(input)
		if len(fields) != 2 {
			fmt.Println("Invalid input.")
			continue
		}

		card1, err := strconv.Atoi(fields[0])
		card1--
		if err != nil || card1 < 0 || card1 >= len(p.Hand) {
			fmt.Println("Invalid input.")
			continue
		}
		firstCard := dealtHand[card1]

		card2, err := strconv.Atoi(fields[1])
		card2--
		if err != nil || card2 == card1 || card2 < 0 || card2 >= len(p.Hand) {
			fmt.Println("Invalid input.")
			continue
		}
		secondCard := dealtHand[card2]

		fmt.Printf("Discarding %s and %s\n", firstCard, secondCard)
		var reset bool
		for {
			reset = false
			fmt.Printf("Continue? [y/n]: ")
			fmt.Scanln(&input)
			switch strings.ToLower(input) {
			case "yes", "y":
				fmt.Println()
				goto exit
			case "no", "n":
				reset = true
				goto exit
			default:
				// repeat loop on "Continue?: "
				continue
			}
		}

	exit:
		if reset {
			PromptIndices(dealtHand)
			// repeat loop on 2-Card selection
			continue
		}

		p.Hand = slices.Delete(p.Hand, card1, card1+1)
		// update index to be the correct element after the array resize
		if card1 < card2 {
			card2--
		}
		p.Hand = slices.Delete(p.Hand, card2, card2+1)
		p.PegHand = make(Hand, len(p.Hand))
		copy(p.PegHand, p.Hand)

		keep = p.PegHand
		discard = Hand{firstCard, secondCard}
		return
	}
}

func (p *HumanPlayer) EmptyPegHand() bool {
	return len(p.PegHand) == 0
}

func (p *HumanPlayer) PlayPegCard(state PegState) (Card, bool) {
	fmt.Println()

	possible := false
	// impossible if PegHand is empty or all cards Ranks are too high
	for i, card := range p.PegHand {
		value := card.ValueMax10()
		if value <= 31-state.Sum {
			possible = true
		}
		fmt.Printf("[%d] %s (value %d)\n", i+1, card, value)
	}
	if !possible {
		// do not play GO automatically but inform the player
		fmt.Println("(You must say Go)")
	}

	fmt.Println("\nSay 'g' to say Go, or 'h' for a Hint on optimal play")
	var input string
	for {
		fmt.Print("Select an index to play that Card: ")
		fmt.Scanln(&input)

		switch strings.ToLower(input) {
		case "go", "g":
			if possible {
				fmt.Println("You have at least 1 valid card and must play!")
				continue
			}
			// blank card and Go/passed is true
			return Card{}, true

		case "help", "h":
			if !possible {
				fmt.Println("You must say Go.")
				continue
			}

			best, ok := OptimalPegging(state, p.PegHand)
			if ok {
				fmt.Printf("The best card to play is %s", best)
				val, _ := ScorePeggingPlay(state, best)
				fmt.Printf(" (+%d)\n", val)
			} else {
				fmt.Println("No optimal play detected.")
			}
			continue
		}

		// index of PegHand (array sized 0-4)
		i, err := strconv.Atoi(input)
		i--
		if err != nil || i < 0 || i >= len(p.PegHand) {
			fmt.Println("Invalid input.")
			// repeat loop on PegHand selection or Say Go
			continue
		}

		card := p.PegHand[i]
		if card.ValueMax10() > 31-state.Sum {
			fmt.Println("Invalid: pile would exceed 31.")
			// repeat loop on PegHand selection or Say Go
			continue
		}

		p.PegHand = slices.Delete(p.PegHand, i, i+1)
		return card, false
	}
}

func (p *HumanPlayer) DrawCard() int {
	fmt.Print("Please select a Card from 1 to 52: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	randomChoice := rand.IntN(52)

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > 52 {
		fmt.Printf("Invalid input: %d was randomly chosen for you\n", randomChoice)
		return randomChoice
	} else {
		return choice - 1
	}
}

func (p *HumanPlayer) CountHand(cut Card, isCrib bool) int {
	reader := bufio.NewReader(os.Stdin)
	var input string
	str_15 := " "
	str_pair := " "
	str_run := " "
	str_flush := " "
	str_nobs := " "

	for input != "d" {
		ClearScreen()
		msg := " Count the number for each Point category "
		if isCrib {
			msg = " (Crib)" + msg
		}
		fmt.Println(msg)
		fmt.Println(strings.Repeat("-", len(msg)))
		fmt.Printf("Hand: %s  Cut: %s\n", p.Hand, cut)

		fmt.Printf("[%s] Fifteens\n[%s] Pairs\n", str_15, str_pair)
		fmt.Printf("[%s] Runs\n[%s] Flush\n[%s] Nobs\n\n", str_run, str_flush, str_nobs)

		fmt.Print("Enter choice (1-5), or 'd' when done: ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > 5 {
			// cannot convert, or input was 'd' for Done
			continue
		}

		switch choice {
		case 1:
			fmt.Print("Number of fifteens [0-8] : ")
			input, _ = reader.ReadString('\n')
			num, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil || num < 0 || num > 8 {
				str_15 = " "
			} else {
				str_15 = fmt.Sprintf("%d", num)
			}

		case 2:
			fmt.Print("Number of pairs [0-6] : ")
			input, _ = reader.ReadString('\n')
			num, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil || num < 0 || num > 6 {
				str_pair = " "
			} else {
				str_pair = fmt.Sprintf("%d", num)
			}

		case 3:
			fmt.Print("Number of runs [0-4] : ")
			input, _ = reader.ReadString('\n')
			num, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil || num < 0 || num > 4 {
				str_run = " "
			} else {
				str_run = fmt.Sprintf("%d", num)
			}

		case 4:
			fmt.Print("Flush [0, 4, 5] : ")
			input, _ = reader.ReadString('\n')
			num, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil || (num != 0 && num != 4 && num != 5) {
				str_flush = " "
			} else {
				str_flush = fmt.Sprintf("%d", num)
			}

		case 5:
			fmt.Print("Nobs [0, 1] : ")
			input, _ = reader.ReadString('\n')
			num, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil || (num != 0 && num != 1) {
				str_nobs = " "
			} else {
				str_nobs = fmt.Sprintf("%d", num)
			}

		default:
			continue
		}
	}

	realPoints := p.Hand.ScoreBreakdown(cut, isCrib)
	var userPoints ScoreBreakdown

	// empty points categories will convert to (0, err)
	fifteenCount, _ := strconv.Atoi(str_15)
	pairCount, _ := strconv.Atoi(str_pair)
	flushPoints, _ := strconv.Atoi(str_flush)
	nobsPoints, _ := strconv.Atoi(str_nobs)
	runCount, _ := strconv.Atoi(str_run)

	userPoints.Fifteens = 2 * fifteenCount
	userPoints.Pairs = 2 * pairCount
	userPoints.Flush = flushPoints
	userPoints.Nobs = nobsPoints
	// save as count because we cannot distinguish Run(s) of 3/4 yet
	userPoints.Runs = runCount

	countedCorrect := CompareBreakDown(userPoints, realPoints)
	if countedCorrect {
		fmt.Println("You counted all points correctly! (NO MUGGINS)")
	}
	fmt.Println()

	if isCrib {
		fmt.Printf("%s (Crib): %s", p.Name, p.Hand)
	} else {
		fmt.Printf("%s: %s", p.Name, p.Hand)
	}
	fmt.Printf(" (%d points)\n", realPoints.Total)
	realPoints.Print()

	fmt.Println()
	p.EnterToContinue()
	return realPoints.Total
}

// equality of two ScoreBreakdown structs, with console messages
func CompareBreakDown(userPoints, realPoints ScoreBreakdown) bool {
	countedCorrect := true

	if realPoints.Fifteens != userPoints.Fifteens {
		countedCorrect = false
		fmt.Printf("Your hand had %d fifteen(s)\n", realPoints.Fifteens/2)
	}

	if realPoints.Pairs != userPoints.Pairs {
		countedCorrect = false
		fmt.Printf("Your hand had %d pair(s)\n", realPoints.Pairs/2)
	}

	if realPoints.Flush != userPoints.Flush {
		countedCorrect = false
		fmt.Printf("Your hand had a flush of %d\n", realPoints.Flush)
	}

	if realPoints.Nobs != userPoints.Nobs {
		countedCorrect = false
		if realPoints.Nobs == 0 {
			fmt.Print("Your hand did not score Nobs")
		} else {
			fmt.Print("Your hand had 1 for His Nobs")
		}
		fmt.Println(" (Jack with a suit matching the Cut Card)")
	}

	runPoints := realPoints.Runs
	// field currently has just the NUMBER of runs input by
	// the User, not actual points (UI does not distinguish run of 3/4/5)
	switch userPoints.Runs {
	case 0:
		if runPoints > 0 {
			countedCorrect = false
			fmt.Println("Uncounted runs")
		} else {
			userPoints.Runs = 0
		}
	case 1:
		// one run of 3 or 4 or 5
		if runPoints != 3 && runPoints != 4 && runPoints != 5 {
			countedCorrect = false
			fmt.Println("Not 1 run")
		} else {
			userPoints.Runs = realPoints.Runs
		}
	case 2:
		// two runs of 3 or 4
		if runPoints != 6 && runPoints != 8 {
			countedCorrect = false
			fmt.Println("Not 2 runs")
		} else {
			userPoints.Runs = realPoints.Runs
		}
	case 3:
		// three runs of 3
		if runPoints != 9 {
			countedCorrect = false
			fmt.Println("Not 3 runs")
		} else {
			userPoints.Runs = 9
		}
	case 4:
		// four runs of 3
		if runPoints != 12 {
			countedCorrect = false
			fmt.Println("Not 4 runs")
		} else {
			userPoints.Runs = 12
		}
	}

	return countedCorrect
}

func (p *HumanPlayer) EnterToContinue() {
	fmt.Print("\nPress any key to continue")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
