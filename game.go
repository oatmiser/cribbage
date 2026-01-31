package cribbage

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func PromptIndices(hand Hand) {
	fmt.Println()
	for i := range hand {
		fmt.Printf(" %d   ", i+1)
	}
	fmt.Println()
	for _, card := range hand {
		fmt.Printf("[%s] ", card)
	}
	fmt.Println()
}

type Player interface {
	// Player selects 4 Hand and 2 Crib cards, mutating state
	Discard(isDealer bool) (Hand, Hand)
	// Player places 1 card from PegHand, mutating state
	PlayPegCard(state PegState) (Card, bool)
	EmptyPegHand() bool
	String() string
	GetName() string
	GetHand() Hand
	SetHand(h Hand)
	AddPoints(n int) int
	GetScore() int
	// select index from shuffled deck of 52
	DrawCard() int
	CountHand(cut Card, isCrib bool) int
	// Human player acknowledges command line outputs
	EnterToContinue()
}

type Game struct {
	Deck    Deck
	Players [2]Player
	Dealer  int
	GameWon bool
}

func (g *Game) ChooseDealer() {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	g.Deck.Shuffle(rng)

	select0 := g.Players[0].DrawCard()
	card0 := g.Deck[select0]
	fmt.Printf("%s pulled %s\n", g.Players[0], card0)

	select1 := g.Players[1].DrawCard()
	if select0 == select1 {
		select1 = (select1 + 1) % 52
	}
	card1 := g.Deck[select1]
	fmt.Printf("%s pulled %s\n", g.Players[1], card1)

	switch {
	case card0.Value() == card1.Value():
		fmt.Println("Cards have the same Rank! Try again...")
		time.Sleep(1 * time.Second)
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
	fmt.Printf("%s will be the first Dealer\n", g.Players[g.Dealer])
	g.Players[0].EnterToContinue()
	g.Players[1].EnterToContinue()
}

func (g *Game) StartGame() {
	roundNum := 0
	for !g.GameWon {
		ClearScreen()
		fmt.Printf("--- Round #%d ---\n", roundNum+1)
		g.PlayRound()
		g.Dealer = 1 - g.Dealer
		roundNum++
	}
}

// add points to player at index i
// Bool indicates "do not continue the Cribbage game" (unused)
func (game *Game) AddPoints(i, amount int) bool {
	total := game.Players[i].AddPoints(amount)
	if total >= 121 {
		game.GameWon = true
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
		// Display player name, points increased, total points
		fmt.Printf("%s (+%d)", player, currentPoints-previous[i])
		fmt.Printf(" : %d points\n", currentPoints)
	}
	fmt.Printf("%s\n", strings.Repeat("-", len(msg)))
}

func (game *Game) CelebrateWinner(winner int) {
	Loser := game.Players[1-winner]
	Winner := game.Players[winner]
	diff := Winner.GetScore() - Loser.GetScore()

	msg := fmt.Sprintf("--- %s won ---", Winner)
	fmt.Printf("\n%s\n", msg)
	for i, player := range game.Players {
		currentPoints := game.Players[i].GetScore()
		fmt.Printf("%s: %d points\n", player, currentPoints)
	}
	fmt.Println(strings.Repeat("-", len(msg)))

	if diff > 60 {
		fmt.Println("DOUBLE SKUNK")
	} else if diff > 30 {
		fmt.Println("SKUNK")
	}
	fmt.Printf("Good Game %s!\n", Loser)
}

func (game *Game) PlayRound() {
	// Each Cribbage round has a shuffle and deal 6 cards
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed)) // pointer
	game.Deck.Shuffle(rng)

	hand1, hand2, remainingDeck := Deal(game.Deck, 6)
	game.Players[0].SetHand(hand1)
	game.Players[1].SetHand(hand2)

	crib := Hand{}
	pone := 1 - game.Dealer
	dealer := game.Dealer

	// Discard to form Crib
	for i, player := range game.Players {
		isDealer := i == dealer
		discard, _ := player.Discard(isDealer)
		crib = append(crib, discard...)

	}

	// Start Pegging round, show Cut card from top of shuffled deck
	ClearScreen()
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

	// Players put one card at a time onto the pile and score points.
	// Make a new pile after cards add to 31 points, until both Hands are empty.
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
	game.Players[0].EnterToContinue()
	game.Players[1].EnterToContinue()
	ClearScreen()

	// Score Hands: Pone will count points first
	fmt.Printf("--- COUNTING ---\n")
	fmt.Printf("Cut Card: %s\n", cut)
	before0 = game.Players[0].GetScore()
	before1 = game.Players[1].GetScore()

	// Score pone's hand
	ponePlayer := game.Players[pone]
	// TODO server calculates points instead
	ponePoints := ponePlayer.CountHand(cut, false)
	// Dealer Player acknowledges the points counted from Pone Computer
	game.Players[dealer].EnterToContinue()
	game.AddPoints(pone, ponePoints)
	if game.GameWon {
		game.CelebrateWinner(pone)
		return
	}

	// Score dealer's hand
	dealPlayer := game.Players[dealer]
	fmt.Printf("Cut Card: %s\n", cut)
	dealerPoints := dealPlayer.CountHand(cut, false)
	game.AddPoints(dealer, dealerPoints)
	if game.GameWon {
		game.CelebrateWinner(dealer)
		return
	}

	// Score dealer's crib, by overwriting their Hand and counting again
	dealPlayer.SetHand(crib)
	cribPoints := dealPlayer.CountHand(cut, true)
	game.AddPoints(dealer, cribPoints)
	if game.GameWon {
		game.CelebrateWinner(dealer)
		return
	}

	fmt.Println()
	game.PrintPoints("SUMMARY (SHOW)", before0, before1)
	game.Players[0].EnterToContinue()
	game.Players[1].EnterToContinue()
	fmt.Println()
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
