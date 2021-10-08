package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Card struct {
	rank string
	suit string
}

type Deck struct {
	Cards []Card
}

type PlayerHand struct {
	Cards []Card
}

type Player struct {
	name        string
	isBusted    bool
	isBlackjack bool
	isStand     bool
	points      int
	PlayerHand
}

type Players []*Player

// calculate the sum of points of the cards
func sumOfPoints(p *Player) int {
	rankToPoints := map[string]int{
		"Ace":   11,
		"Two":   2,
		"Three": 3,
		"Four":  4,
		"Five":  5,
		"Six":   6,
		"Seven": 7,
		"Eight": 8,
		"Nine":  9,
		"Ten":   10,
		"Jack":  10,
		"Queen": 10,
		"King":  10,
	}

	// initial calculation without taking account of the Ace rules
	totalPoints := 0
	func() {
		for _, c := range p.Cards {
			totalPoints += rankToPoints[c.rank]
		}
	}()

	// Ace is 1 or 10 points when the player has 3 cards - depending on whether the player is busted
	// Ace is 1 point when player has more than 3 cards
	switch {
	case len(p.Cards) == 3: // if player has 3 cards
		rankToPoints["Ace"] = 10 // if player is not busted then Ace is 10 points
		if p.points > 21 {       // if player is busted then Ace becomes 1 point
			rankToPoints["Ace"] = 1
		}
	case len(p.Cards) > 3: // Ace is 1 point when player has more than 3 cards
		rankToPoints["Ace"] = 1
	}

	// recalculate using the newly assigned Ace points
	totalPoints = 0
	func() {
		for _, c := range p.Cards {
			totalPoints += rankToPoints[c.rank]
		}
	}()

	return totalPoints
}

// deal a card to a player
func dealCard(d *Deck, p *Player) {
	p.Cards = append(p.Cards, d.Cards[0]) // add the deck's top card to the player's hand
	d.Cards = d.Cards[1:]                 // remove the top card from the deck
	updatePlayerState(p)
	reportAfterDealCard((p))
}

// updates the states of the player
func updatePlayerState(p *Player) {
	p.points = sumOfPoints(p)
	p.isBlackjack = isBlackjack(p)
	p.isBusted = isBusted(p)
}

// returns a fresh deck of cards
func newDeck() Deck {
	suits := []string{"Clubs", "Diamonds", "Hearts", "Spades"}
	ranks := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
	deck := Deck{}
	func() {
		for _, s := range suits {
			for _, r := range ranks {
				var newCard Card
				newCard.rank = r
				newCard.suit = s
				deck.Cards = append(deck.Cards, newCard)
			}
		}
	}()
	return deck
}

// shuffle the deck of cards
func shuffle(d *Deck) {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// decide if a player's hand is busted
func isBusted(p *Player) bool {
	return sumOfPoints(p) > 21
}

// decide if a player's hand is a blackjack
func isBlackjack(p *Player) bool {
	if len(p.Cards) == 2 && sumOfPoints(p) == 21 {
		return true
	}
	return false
}

func reportDealtCard(p *Player) {
	mostRecentCard := p.Cards[len(p.Cards)-1]
	fmt.Printf("[%s] - Card [%v] - [%s] of [%s].\n", p.name, len(p.Cards), mostRecentCard.rank, mostRecentCard.suit)
}

func reportPlayerState(p *Player) {
	if len(p.Cards) == 2 {
		fmt.Printf("[%s] - Total points: [%v]. Blackjack: [%v].\n", p.name, p.points, p.isBlackjack)
	} else if len(p.Cards) > 2 {
		fmt.Printf("[%s] - Total points: [%v]. Busted: [%v].\n", p.name, p.points, p.isBusted)
	}
}

func reportAfterDealCard(p *Player) {
	reportDealtCard(p)
	reportPlayerState(p)
}

// ask whether the players want to hit or stand
func hitOrStand(d *Deck, p *Player) string {
	aceCount := 0
	for _, c := range p.Cards {
		if c.rank == "Ace" {
			aceCount++
		}
	}
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("[%s] - Current Points: [%v], Card Count: [%v], Number of Aces: [%v]\n", p.name, p.points, len(p.Cards), aceCount)
		if p.points < 16 {
			fmt.Printf("[%s] - Not enough points. Must Hit(H): ", p.name)
		} else {
			fmt.Printf("[%s] - Hit(H) or Stand(S): ", p.name)
		}
		scanner.Scan()
		decision := scanner.Text()
		if !(decision == "H" || decision == "S") {
			fmt.Print("Invalid input detected. You need to type 'H' or 'S'!\n")
		} else {
			return decision
		}
	}
}

func startGame(d *Deck, players Players) {
	// deal two cards to each player at the start
	for _, p := range players {
		dealCard(d, p)
	}
	for _, p := range players {
		dealCard(d, p)
	}

	for _, p := range players {
	out:
		// max number of cards to be dealt to any player is 5
		for len(p.Cards) <= 5 {
			switch {
			case len(p.Cards) >= 2 && len(p.Cards) <= 5 && playerHasMoves(p): // third, fourth and fifth card
				decision := hitOrStand(d, p) // ask whether the player wants to hit or stand
				if decision == "H" {
					dealCard(d, p)
				} else if decision == "S" {
					p.isStand = true
				}
			case !playerHasMoves(p):
				break out // break to outer loop
			}
		}
		fmt.Printf("[%s] has no more moves. Moving to next player.\n", p.name)
	}
	fmt.Println("Calculating the results...")
	calculateOutcome(players)
}

// get the winner by comparing the player and the dealer
func getWinner(player *Player, dealer *Player) string {
	if player.isBlackjack && !dealer.isBlackjack {
		return player.name
	} else if !player.isBusted && player.points > dealer.points {
		return player.name
	} else if !player.isBusted && dealer.isBusted {
		return player.name
	} else if player.isBlackjack && dealer.isBlackjack {
		return ""
	} else if player.isBusted && dealer.isBusted {
		return ""
	} else if player.points == dealer.points {
		return ""
	}
	return dealer.name
}

// calculate the outcome of the game
func calculateOutcome(players Players) {
	for _, p := range players {
		if p.name != "Dealer" {
			winner := getWinner(p, players[len(players)-1])
			if winner != "" {
				fmt.Printf("[%s]: %v vs Dealer: %v - Winner: %s.\n", p.name, p.points, players[len(players)-1].points, winner)
			} else {
				fmt.Printf("[%s]: %v vs Dealer: %v - Draw.\n", p.name, p.points, players[len(players)-1].points)
			}
		}
	}
}

// return whether there is any remaining moves that the player could possibly take
func playerHasMoves(p *Player) bool {
	return !(p.isBlackjack || p.isBusted || p.isStand)
}

func main() {
	rand.Seed(time.Now().UnixNano()) // set a random seed based on current time
	fmt.Println("Getting a new deck...")
	deck := newDeck()

	fmt.Println("Shuffling deck...")
	shuffle(&deck)

	players := Players{}
	dealer := Player{name: "Dealer"}
	playerA := Player{name: "Alan"}
	playerB := Player{name: "Bob"}

	players = append(players, &playerA, &playerB, &dealer) // dealer goes last

	startGame(&deck, players)
}
