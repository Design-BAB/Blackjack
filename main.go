//Author: Design-BAB
//Date: 1/2/2025
//Description: Time to play Blackjack! The goal is to reach 452 lines of code

package main

import (
	"math/rand"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowSize        = 800
	Center            = WindowSize / 2
	PlayerCardXPos    = 175
	PlayerCardYPos    = WindowSize - 200
	MaxFrames         = 432000
	TotalDeck         = 52
	MaxHand           = 11
	Blackjack         = 21
	TimeToPassOutCard = 450 * time.Millisecond
)

type GameState struct {
	JustStarted        bool
	IsOver             bool
	Lives              int
	TurnIsNow          bool
	RoundIsOver        bool
	YouWon             bool
	DealerGettingCards bool
	Scheduler          time.Time
	LastCardTime       time.Time
}

func newGame() *GameState {
	startTimeNow := time.Now()
	return &GameState{JustStarted: true, Lives: 3, Scheduler: startTimeNow, LastCardTime: startTimeNow}
}

type Player struct {
	Hand     [MaxHand]*CardInHand
	Score    int
	IsDealer bool
	IsBust   bool
}

func newPlayer(hand [MaxHand]*CardInHand, isDealer bool) *Player {
	return &Player{Hand: hand, IsDealer: isDealer}
}

func calculateScore(player *Player) {
	theScore := 0
	numberOfAces := 0
	//Calculate score with ace being default 11
	for i := range MaxHand {
		if player.Hand[i] == nil {
			break
		} else {
			theScore += player.Hand[i].Value
			if player.Hand[i].Value == 11 {
				numberOfAces++
			}
		}
	}
	//if they go over place jack, then we will have the ace go being valued 1
	if numberOfAces > 0 {
		for i := 0; i < numberOfAces; i++ {
			if theScore > Blackjack {
				theScore = theScore - 10
			}
		}
	}
	player.Score = theScore
}

type Card struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	IsDiscarded  bool
	Value        int
	BelongsToYou bool
}

func newCard(texture rl.Texture2D, x, y float32, value int) *Card {
	return &Card{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Value: value, BelongsToYou: true}
}

type CardInHand struct {
	*Card
	ToShow bool
}

func AddCardToHand(theCard *Card) *CardInHand {
	return &CardInHand{Card: theCard}
}

// this function will run in main in order to import cards
func importCards() ([TotalDeck]rl.Texture2D, [TotalDeck]*Card) {
	var cardTextures [TotalDeck]rl.Texture2D
	var cardDeck [TotalDeck]*Card
	c := 0
	cardTextures, cardDeck, c = loadCardTexture("Clubs", c, cardTextures, cardDeck)
	cardTextures, cardDeck, c = loadCardTexture("Diamonds", c, cardTextures, cardDeck)
	cardTextures, cardDeck, c = loadCardTexture("Hearts", c, cardTextures, cardDeck)
	cardTextures, cardDeck, c = loadCardTexture("Spades", c, cardTextures, cardDeck)
	return cardTextures, cardDeck
}

func loadCardTexture(typeOfCard string, c int, cardTextures [TotalDeck]rl.Texture2D, cardDeck [TotalDeck]*Card) ([TotalDeck]rl.Texture2D, [TotalDeck]*Card, int) {
	//could this function be done a lot shorter? probably but I'm doing 2 switch statements because sanity and 2 is a number I can live with
	for i := 2; i < 15; i++ {
		//first switch statement for the textures
		switch i {
		case 11:
			cardTextures[c] = rl.LoadTexture("images/Jack_" + typeOfCard + ".png")
		case 12:
			cardTextures[c] = rl.LoadTexture("images/Queen_" + typeOfCard + ".png")
		case 13:
			cardTextures[c] = rl.LoadTexture("images/King_" + typeOfCard + ".png")
		case 14:
			cardTextures[c] = rl.LoadTexture("images/Ace_" + typeOfCard + ".png")
		default:
			cardTextures[c] = rl.LoadTexture("images/" + strconv.Itoa(i) + "_" + typeOfCard + ".png")
		}
		//second switch for the deck struct
		switch i {
		case 11, 12, 13:
			cardDeck[c] = newCard(cardTextures[c], PlayerCardXPos, PlayerCardYPos, 10)
		case 14:
			cardDeck[c] = newCard(cardTextures[c], PlayerCardXPos, PlayerCardYPos, 11)
		default:
			cardDeck[c] = newCard(cardTextures[c], PlayerCardXPos, PlayerCardYPos, i)
		}
		//ha ha ha no, this is not a programming pun. c refers to card count therefore c++
		c++
	}
	return cardTextures, cardDeck, c
}

// this grabs user's input. Keep in mind that the scheduler is set off here
func getInput(player1, dealer *Player, cardDeck [TotalDeck]*Card, yourGame *GameState) {
	if yourGame.IsOver == false && yourGame.JustStarted == false && yourGame.TurnIsNow {
		if rl.IsKeyPressed(rl.KeyX) {
			hit(player1, cardDeck)
			if player1.Score > Blackjack {
				// Player busted!
				yourGame.TurnIsNow = false
				yourGame.RoundIsOver = true
				yourGame.Lives--
				yourGame.Scheduler = time.Now()
			}
		}
		if rl.IsKeyPressed(rl.KeyZ) {
			stay(player1, dealer, cardDeck, yourGame)
		}
	} else if yourGame.IsOver {
		if rl.IsKeyPressed(rl.KeyA) {
			yourGame.IsOver = false
			yourGame.Lives = 3
			yourGame.JustStarted = true
			yourGame.RoundIsOver = false
			yourGame.TurnIsNow = true
			player1.Score = 0
			dealer.Score = 0
			resetRound(player1, dealer, yourGame)
		}
	}
}

func hit(player *Player, cardDeck [TotalDeck]*Card) {
	var newCardi int
	//first find the empty slot
	for i := range MaxHand {
		if player.Hand[i] == nil {
			newCardi = i
			break
		}
	}
	//now we grab a new card
	for i := range TotalDeck {
		if cardDeck[i].IsDiscarded == false {
			player.Hand[newCardi] = AddCardToHand(cardDeck[i])
			player.Hand[newCardi].ToShow = true
			cardDeck[i].IsDiscarded = true
			//change the position if it is the dealer's cards
			if player.IsDealer {
				player.Hand[newCardi].X = WindowSize - 5 - player.Hand[newCardi].Width
				player.Hand[newCardi].Y = 75
			}
			break
		}
	}
	calculateScore(player)
}

func stay(player1, dealer *Player, cardDeck [TotalDeck]*Card, yourGame *GameState) {
	yourGame.TurnIsNow = false
	yourGame.DealerGettingCards = true
	yourGame.LastCardTime = time.Now()
	calculateScore(player1)
	calculateScore(dealer)
}

func update(cardDeck [TotalDeck]*Card, player1, dealer *Player, yourGame *GameState) {
	if yourGame.IsOver == false {
		if yourGame.JustStarted {
			rand.Shuffle(TotalDeck, func(i, j int) {
				cardDeck[i], cardDeck[j] = cardDeck[j], cardDeck[i]
			})
			//gotta reset the deck if we run out of cards.
			if cardDeck[TotalDeck-10].IsDiscarded == true {
				for i := range TotalDeck {
					cardDeck[i].IsDiscarded = false
					cardDeck[i].X = PlayerCardXPos
					cardDeck[i].Y = PlayerCardYPos
				}
			}
			//giving first card and second card
			hit(player1, cardDeck)
			hit(player1, cardDeck)
			//giving first & second card to dealer
			hit(dealer, cardDeck)
			hit(dealer, cardDeck)
			//update Score
			calculateScore(player1)
			calculateScore(dealer)
			yourGame.JustStarted = false
			yourGame.TurnIsNow = true
		}
		if yourGame.DealerGettingCards {
			if time.Since(yourGame.LastCardTime) > TimeToPassOutCard {
				if dealer.Score < 17 {
					hit(dealer, cardDeck)
					yourGame.LastCardTime = time.Now()
				} else {
					yourGame.DealerGettingCards = false
					checkWinConditions(player1, dealer, yourGame)
					yourGame.RoundIsOver = true
				}
			}
		}
		//remember: the checking of win conditions happen in checkWinConditions while
		//losing conidtions is checked in that and getInput
		if yourGame.Lives <= 0 {
			yourGame.IsOver = true
			yourGame.RoundIsOver = false
		}
	}
}

func resetRound(player1, dealer *Player, yourGame *GameState) {
	for i := range MaxHand {
		if player1.Hand[i] != nil {
			player1.Hand[i] = nil
		}
		if dealer.Hand[i] != nil {
			dealer.Hand[i] = nil
		}
	}
	yourGame.JustStarted = true
	yourGame.TurnIsNow = true
}

func checkWinConditions(player1, dealer *Player, yourGame *GameState) {
	if player1.Score == Blackjack && dealer.Score != Blackjack {
		yourGame.YouWon = true
	} else if dealer.Score > Blackjack && player1.Score < Blackjack {
		yourGame.YouWon = true
	} else if dealer.Score < player1.Score && dealer.Score < Blackjack && player1.Score < Blackjack {
		yourGame.YouWon = true
	} else {
		//if they didn't meet the win conditions, then it is safe to say they lost
		yourGame.Lives--
	}
	startTimeNow := time.Now()
	yourGame.Scheduler = startTimeNow
}

func draw(background, backOfCard, heart rl.Texture2D, player1, dealer *Player, yourGame *GameState) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	//rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	rl.DrawTexture(background, 0, 0, rl.White)
	if yourGame.IsOver {
		rl.DrawText("Game is over", 190, 200, 20, rl.Red)
		rl.DrawText("Press 'A' to reset", 190, 250, 20, rl.Red)

	} else {
		offset := int32(0)
		for i := range MaxHand {
			if player1.Hand[i] == nil {
				break
			} else {
				if player1.Hand[i] != nil && player1.Hand[i].ToShow == true {
					rl.DrawTexture(player1.Hand[i].Texture, int32(player1.Hand[i].X)+offset, int32(player1.Hand[i].Y), rl.White)
					offset = offset + 30
				}
			}
		}
		rl.DrawText("Score: "+strconv.Itoa(player1.Score), 10, WindowSize-75, 25, rl.RayWhite)
		offset = 0
		isFirst := true
		for i := range MaxHand {
			if dealer.Hand[i] == nil {
				break
			} else {
				if dealer.Hand[i] != nil && dealer.Hand[i].ToShow == true {
					if isFirst && yourGame.RoundIsOver == false {
						rl.DrawTexture(backOfCard, int32(dealer.Hand[i].X), int32(dealer.Hand[i].Y), rl.White)
						offset = offset + 40
						isFirst = false
					} else {
						rl.DrawTexture(dealer.Hand[i].Texture, int32(dealer.Hand[i].X), int32(dealer.Hand[i].Y)+offset, rl.White)
						offset = offset + 40
					}
				}
			}
		}
		if yourGame.RoundIsOver {
			rl.DrawText(strconv.Itoa(dealer.Score), WindowSize-100, 5, 20, rl.RayWhite)
			if time.Since(yourGame.Scheduler) < 3*time.Second {
				if yourGame.YouWon {
					rl.DrawText("You won!", 25, 180, 40, rl.Green)
				} else if dealer.Score == player1.Score {
					rl.DrawText("The house wins", 25, 180, 35, rl.Red)
				} else {
					rl.DrawText("You lost", 25, 180, 40, rl.Red)
				}
			} else {
				yourGame.YouWon = false
				yourGame.RoundIsOver = false
				yourGame.TurnIsNow = true
				resetRound(player1, dealer, yourGame)
			}
		}
	}
	drawUi(heart, yourGame)
	rl.EndDrawing()
}

func drawUi(heart rl.Texture2D, yourGame *GameState) {
	rl.DrawTexture(heart, 60, WindowSize-30, rl.White)
	rl.DrawText(strconv.Itoa(yourGame.Lives), 90, WindowSize-30, 40, rl.RayWhite)
	//The variable yourGame.JustStarted is tracking when to hand out 2 cards but turns off immediatly
	//to prevent the loop to go crazy... so we are using lives to detect when the begining of the game is
	if yourGame.Lives == 3 {
		rl.DrawText("Press 'X' to hit", 5, 5, 20, rl.RayWhite)
		rl.DrawText("Press Z to stay", 5, 30, 20, rl.RayWhite)

	}
}

func main() {
	rl.InitWindow(WindowSize, WindowSize, "Blackjack")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	yourGame := newGame()
	//import textures
	background := rl.LoadTexture("images/background.png")
	defer rl.UnloadTexture(background)
	backOfCard := rl.LoadTexture("images/Backface_Red.png")
	defer rl.UnloadTexture(backOfCard)
	heart := rl.LoadTexture("images/life-count.png")
	defer rl.UnloadTexture(heart)
	//importing cards
	cardTextures, cardDeck := importCards()
	for _, texture := range cardTextures {
		defer rl.UnloadTexture(texture)
	}
	//setting up hands and players
	var yourHand [MaxHand]*CardInHand
	var herHand [MaxHand]*CardInHand
	player1 := newPlayer(yourHand, false)
	dealer := newPlayer(herHand, true)
	frames := 0
	for !rl.WindowShouldClose() && frames < MaxFrames {
		if yourGame.RoundIsOver == false {
			getInput(player1, dealer, cardDeck, yourGame)
		}
		update(cardDeck, player1, dealer, yourGame)
		draw(background, backOfCard, heart, player1, dealer, yourGame)
		frames++
	}
}
