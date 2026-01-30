package cribbage

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Player interface {
	// Player selects 4 Hand and 2 Crib cards, mutating state
	Discard(isDealer bool) (Hand, Hand)
	PlayPegCard(state PegState) (Card, bool)
	String() string
	GetName() string
	GetHand() Hand
	SetHand(h Hand)
	AddPoints(n int) int
	GetScore() int
	DrawCard() int // select index from shuffled deck of 52
	CountHand(cut Card, isCrib bool) int
	EnterToContinue()
	EmptyPegHand() bool
}

type Game struct {
	Deck    Deck
	Players [2]Player
	Dealer  int
	GameWon bool
}

func (g *Game) ChooseDealer() {
	//ClearScreen()
	g.Deck.Shuffle()

	select0 := g.Players[0].DrawCard()
	name0 := g.Players[0].GetName()
	card0 := g.Deck[select0]
	//fmt.Printf("%s pulled %s (index %d)\n", name0, card0, select0+1)
	fmt.Printf("%s pulled %s\n", name0, card0)

	select1 := g.Players[1].DrawCard()
	// player goes first in a PlayerGame so we don't need the message
	if select0 == select1 {
		//fmt.Println("Cannot select the same card! (defaulting to next index)")
		select1 = (select1 + 1) % 52
		//fmt.Printf("%s already chose %s; ", name0, g.Deck[select0])
		//fmt.Printf("%s take next card %s\n", name1, g.Deck[select1])
	}

	name1 := g.Players[1].GetName()
	card1 := g.Deck[select1]
	//fmt.Printf("%s pulled %s (index %d)\n", name1, card1, select1+1)
	fmt.Printf("%s pulled %s\n", name1, card1)

	switch {
	case card0.Value() == card1.Value():
		fmt.Println("Cards have the same Rank! Try again...")
		g.ChooseDealer()
	case card0.Value() < card1.Value():
		g.Dealer = 0
		fmt.Printf("\n%s is the lower rank; ", card0)
	case card0.Value() > card1.Value():
		g.Dealer = 1
		fmt.Printf("\n%s is the lower rank; ", card1)
	default:
		fmt.Println("TODO ChooseDealer")
		g.Dealer = 0
	}
	fmt.Printf("%s will be the first Dealer\n\n", g.Players[g.Dealer].GetName())
	g.Players[0].EnterToContinue()
	g.Players[1].EnterToContinue()
	//time.Sleep(2 * time.Second)
}

func (g *Game) StartGame() {
	roundNum := 0
	//for g.Players[0].GetScore() < 121 && g.Players[1].GetScore() < 121 {
	for !g.GameWon {
		ClearScreen()
		fmt.Printf("--- Round #%d ---\n", roundNum+1)
		g.PlayRound()
		g.Dealer = 1 - g.Dealer
		roundNum++
	}
}

func (game *Game) AddPoints(i, amount int) bool {
	total := game.Players[i].AddPoints(amount)
	if total >= 121 {
		game.GameWon = true
		// do not continue the Cribbage round (unused)
		return false
	}
	return true
}

// Print total points with some message/header
func (g *Game) PrintPoints(msg string, previous0 int, previous1 int) {
	previous := [2]int{previous0, previous1}
	msg = fmt.Sprintf("--- %s ---", msg)

	fmt.Println(msg)
	for i, player := range g.Players {
		currentPoints := g.Players[i].GetScore()
		// Player name, points increased, total points
		fmt.Printf("%s (+%d)", player, currentPoints-previous[i])
		fmt.Printf(" : %d points\n", currentPoints)
	}
	fmt.Printf("%s\n", strings.Repeat("-", len(msg)))
}

func (game *Game) CelebrateWinner(winner int) {
	Loser := game.Players[1-winner]
	Winner := game.Players[winner]
	diff := Winner.GetScore() - Loser.GetScore()

	msg := fmt.Sprintf("--- %s won ---", Winner.GetName())
	fmt.Printf("\n%s\n", msg)
	for i, player := range game.Players {
		currentPoints := game.Players[i].GetScore()
		fmt.Printf("%s: %d points\n", player.GetName(), currentPoints)
	}
	fmt.Println(strings.Repeat("-", len(msg)))

	if diff > 60 {
		fmt.Println("DOUBLE SKUNK")
	} else if diff > 30 {
		fmt.Println("SKUNK")
	}
	fmt.Printf("Good Game %s!\n", Loser.GetName())
}

func (game *Game) PlayRound() {
	// Simulate shuffle and card deal
	game.Deck.Shuffle()
	hand1, hand2, remainingDeck := Deal(game.Deck, 6)
	game.Players[0].SetHand(hand1)
	game.Players[1].SetHand(hand2)
	//ShowCardDeal(game.Players[0], game.Players[1])

	crib := Hand{}
	pone := 1 - game.Dealer
	dealer := game.Dealer

	// Discard to form Crib
	for i, player := range game.Players {
		isDealer := i == dealer
		discard, _ := player.Discard(isDealer)
		//fmt.Println(keep)
		crib = append(crib, discard...)

	}

	// Start Pegging round

	ClearScreen()
	// Cut card from top of deck
	cut := remainingDeck[0]
	fmt.Printf("Cut Card: %s\n", cut)
	if cut.Rank == Jack {
		fmt.Println("TWO FOR HIS HEELS")
		fmt.Printf("%s scores +2 [Nibs]\n", game.Players[dealer])
		game.AddPoints(dealer, 2)
		if game.GameWon {
			game.CelebrateWinner(dealer)
			return
		}
	}

	// Player's place one card at a time onto a pile and score points
	// New pile after cards add to 31 points, until both Hands are empty
	title := "--- PEGGING ---"
	fmt.Printf("\n%s\n", title)
	fmt.Printf("Dealer: %s\nPone leads\n", game.Players[dealer])
	fmt.Printf("%s\n\n", strings.Repeat("-", len(title)))

	before0 := game.Players[0].GetScore()
	before1 := game.Players[1].GetScore()
	game.StartPegging()
	if game.GameWon {
		// CelebrateWinner inside StartPegging
		return
	}
	fmt.Println()
	game.PrintPoints("SUMMARY (PLAY)", before0, before1)
	//time.Sleep(2 * time.Second)
	fmt.Println()
	game.Players[0].EnterToContinue()
	game.Players[1].EnterToContinue()
	ClearScreen()

	// Score Hands: Pone will count points first
	fmt.Printf("--- COUNTING ---\n")
	fmt.Printf("Cut Card: %s\n", cut)
	before0 = game.Players[0].GetScore()
	before1 = game.Players[1].GetScore()
	//time.Sleep(2 * time.Second)

	// Score pone's hand
	ponePlayer := game.Players[pone]
	ponePoints := ponePlayer.CountHand(cut, false)
	/*poneHand := ponePlayer.GetHand()
	ponePoints := poneHand.ScoreBreakdown(cut, false)
	fmt.Printf("%s: %s", ponePlayer.GetName(), poneHand)
	fmt.Printf(" (%d points)\n", ponePoints.Total)
	ponePoints.Print()
	*/

	// Player acknowledges points from Computer when it shows first
	fmt.Println()
	game.Players[dealer].EnterToContinue()

	game.AddPoints(pone, ponePoints) //.Total)
	if game.GameWon {
		game.CelebrateWinner(pone)
		return
	}

	// Score dealer's hand
	dealPlayer := game.Players[dealer]
	fmt.Printf("Cut Card: %s\n", cut)
	dealerPoints := dealPlayer.CountHand(cut, false)
	/*dealerHand := dealPlayer.GetHand()
	dealerPoints := dealerHand.ScoreBreakdown(cut, false)
	fmt.Printf("%s: %s", dealPlayer.GetName(), dealerHand)
	fmt.Printf(" (%d points)\n", dealerPoints.Total)
	dealerPoints.Print()
	*/

	game.AddPoints(dealer, dealerPoints)
	if game.GameWon {
		game.CelebrateWinner(dealer)
		return
	}

	// Score dealer's crib
	// set hand of dealer and count again
	dealPlayer.SetHand(crib)
	cribPoints := dealPlayer.CountHand(cut, true)
	/*cribPoints := crib.ScoreBreakdown(cut, true)
	fmt.Printf("%s (Crib): %s", dealPlayer.GetName(), crib)
	fmt.Printf(" (%d points)\n", cribPoints.Total)
	cribPoints.Print()
	*/

	game.AddPoints(dealer, cribPoints) //.Total)
	if game.GameWon {
		game.CelebrateWinner(dealer)
		return
	}

	fmt.Println()
	game.PrintPoints("SUMMARY (SHOW)", before0, before1)

	fmt.Println()
	game.Players[0].EnterToContinue()
	game.Players[1].EnterToContinue()
	fmt.Println()
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

// ------------------------------------------------------------ //

func NewComputerGame() *Game {
	p1 := &ComputerPlayer{
		Name:   "COM 1",
		Points: 0,
	}

	p2 := &ComputerPlayer{
		Name:   "COM 2",
		Points: 0,
	}

	return &Game{
		Deck:    NewDeck(),
		Players: [2]Player{p1, p2},
		Dealer:  0,
	}
}

func NewPlayerGame() *Game {

	fmt.Print("Please provide your name: ")
	//fmt.Scanln(&name)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	name := strings.TrimSpace(input)

	p1 := &HumanPlayer{
		Name:   name,
		Points: 0,
	}

	p2 := &ComputerPlayer{
		Name:   "COM 1",
		Points: 0,
	}

	return &Game{
		Deck:    NewDeck(),
		Players: [2]Player{p1, p2},
		Dealer:  0,
	}
}

func Start() {

	AwesomeTitle := `
  ___                                         
 / _ \                                        
/ /_\ \_      _____  ___  ___  _ __ ___   ___ 
|  _  \ \ /\ / / _ \/ __|/ _ \| '_ ' _ \ / _ \
| | | |\ V  V /  __/\__ \ (_) | | | | | |  __/
\_| |_/ \_/\_/ \___||___/\___/|_| |_| |_|\___|

 _____      _ _     _                         
/  __ \    (_) |   | |                        
| /  \/_ __ _| |__ | |__   __ _  __ _  ___    
| |   | '__| | '_ \| '_ \ / _' |/ _' |/ _ \   
| \__/\ |  | | |_) | |_) | (_| | (_| |  __/   
 \____/_|  |_|_.__/|_.__/ \__,_|\__, |\___|   
                                 __/ |        
                                |___/         
 _____                                        
|  __ \                                       
| |  \/ __ _ _ __ ___   ___                   
| | __ / _' | '_ ' _ \ / _ \                  
| |_\ \ (_| | | | | | |  __/                  
 \____/\__,_|_| |_| |_|\___|  v1
`
	// Awesome Cribbage Game
	ClearScreen()
	fmt.Println(AwesomeTitle)
	var game *Game
	var input string

	fmt.Print("Play cribbage with input? ")
	fmt.Scanln(&input)
	switch strings.ToLower(input) {
	case "yes", "y":
		game = NewPlayerGame()
	default:
		game = NewComputerGame()
	}

	fmt.Printf("Welcome %s and %s!\n", game.Players[0], game.Players[1])
	fmt.Printf("A new game is beginning...\n\n")
	time.Sleep(1 * time.Second)
	game.ChooseDealer()
	game.StartGame()
}
