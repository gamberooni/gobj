# Go Blackjack

This is my first `Go` project that implements a simplified version of the Blackjack game.

## Rules

- Each Player versus the Dealer individually.
- Maximum number of cards dealt per Player/Dealer is 5.
- All cards are faced up, including the Dealer's.
- Player/Dealer must `Hit` when points is lower than **16**.
- Player/Dealer can decide to `Hit` or `Stand` when points is greater than or equal to **16**.
- Player/Dealer is busted when points is greater than **21**.
- Ace rules:
  - Worths **11** points when the Player/Dealer has **2** cards.
  - Worths **10** points when the Player/Dealer has **3** cards and is **not busted**.
  - Worths **1** point when the Player/Dealer has **3** cards and **should have been busted** if the Ace was counted as 10 points.
  - Worths **1** point when Player/Dealer has **more than 3 cards**.
- No bets.
- No splits.

## Game Flow

1. Pull out a fresh deck of cards
2. Shuffle the deck using a random seed based on the current time.
3. Deal 2 cards to each player. **Dealer goes last**.
4. Ask for decision to `Hit` or `Stand` starting from the first Player until the Player has no remaining moves.
    - Player has no remaining moves when the Player gets a Blackjack, decides to `Stand` or is busted.
5. Proceeds to the next player.
6. Repeat steps (4) to (5) until there are no Players left (Dealer is the last "Player").
7. Calculate the outcome of the game.

## How to Start Game

Clone this repo and run `go run main.go`. By default, there are 2 Players (Alan and Bob) and a Dealer. You can edit the `main` function to add more Players if you wish to. **Note that the Dealer is the last "Player"!**
