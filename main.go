//Author: Design-BAB
//Date: 12/25/2025
//Description: Time to play Blackjack! The goal is to reach 452 lines of code

package main

import (
	"fmt"
	"math/rand"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowSize = 800
	Center     = WindowSize / 2
	MaxFrames  = 432000
	TotalDeck  = 52
	MaxHand    = 11
	Blackjack  = 21
)

type GameState struct {
	JustStarted bool
	IsOver      bool
	Lives       int
}

func newGame() *GameState {
	return &GameState{JustStarted: true, Lives: 3}
}

type Player struct {
	Hand      [MaxHand]*CardInHand
	Score     int
	IsDealer  bool
	IsBust    bool
	HasStayed bool
}

func newPlayer(hand [MaxHand]*CardInHand, isDealer bool) *Player {
	return &Player{Hand: hand, IsDealer: isDealer}
}

func calculateScore(player *Player) {
	theScore := 0
	numberOfAces := 0
	for i := 0; i < MaxHand; i++ {
		if player.Hand[i] == nil {
			break
		} else {
			theScore += player.Hand[i].Value
			if player.Hand[i].Value == 11 {
				numberOfAces++
			}
		}
	}
	if numberOfAces > 0 {
		for i := 0; i < numberOfAces; i++ {
			if theScore > Blackjack {
				theScore = theScore - 10
			}
		}
	}
	fmt.Println(theScore)
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
		//second switch for the deck
		switch i {
		case 11, 12, 13:
			cardDeck[c] = newCard(cardTextures[c], Center-200, WindowSize-200, 10)
		case 14:
			cardDeck[c] = newCard(cardTextures[c], Center-200, WindowSize-200, 11)
		default:
			cardDeck[c] = newCard(cardTextures[c], Center-200, WindowSize-200, i)
		}
		//ha ha ha no, this is not a programming pun. c refers to card count therefore c++
		c++
	}
	return cardTextures, cardDeck, c
}

func getInput(player1 *Player, cardDeck [TotalDeck]*Card, yourGame *GameState) {
	if yourGame.IsOver == false && yourGame.JustStarted == false {
		if rl.IsKeyPressed(rl.KeyX) {
			hit(player1, cardDeck)
		}
	}
}

func hit(player *Player, cardDeck [TotalDeck]*Card) {
	var newCardi int
	//first find the empty slot
	for i := 0; i < MaxHand; i++ {
		if player.Hand[i] == nil {
			newCardi = i
			break
		}
	}
	//now we grab a new card
	for i := 0; i < TotalDeck; i++ {
		if cardDeck[i].IsDiscarded == false {
			player.Hand[newCardi] = AddCardToHand(cardDeck[i])
			player.Hand[newCardi].ToShow = true
			cardDeck[i].IsDiscarded = true
			break
		}
	}
	calculateScore(player)
}

func update(cardDeck [TotalDeck]*Card, player1, dealer *Player, yourGame *GameState) {
	if yourGame.IsOver == false {
		if yourGame.JustStarted {
			rand.Shuffle(TotalDeck, func(i, j int) {
				cardDeck[i], cardDeck[j] = cardDeck[j], cardDeck[i]
			})
			//giving first card
			player1.Hand[0] = AddCardToHand(cardDeck[0])
			player1.Hand[0].ToShow = true
			cardDeck[0].IsDiscarded = true
			//giving second card
			player1.Hand[1] = AddCardToHand(cardDeck[1])
			player1.Hand[1].ToShow = true
			cardDeck[1].IsDiscarded = true
			//giving first card to dealer
			dealer.Hand[0] = AddCardToHand(cardDeck[2])
			dealer.Hand[0].ToShow = true
			dealer.Hand[0].X = WindowSize - 5 - dealer.Hand[0].Width
			dealer.Hand[0].Y = 50
			cardDeck[2].IsDiscarded = true
			//giving second card to dealer
			dealer.Hand[1] = AddCardToHand(cardDeck[3])
			dealer.Hand[1].ToShow = true
			dealer.Hand[1].X = WindowSize - 5 - dealer.Hand[1].Width
			dealer.Hand[1].Y = 50
			cardDeck[3].IsDiscarded = true
			//update Score
			calculateScore(player1)
			calculateScore(dealer)
			yourGame.JustStarted = false
		}
		if player1.Score > Blackjack {
			yourGame.Lives -= 1
		}
		if yourGame.Lives <= 0 {
			yourGame.IsOver = true
		}
	}
}

func draw(background, backOfCard rl.Texture2D, player1, dealer *Player, yourGame *GameState) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	//rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	rl.DrawTexture(background, 0, 0, rl.White)
	if yourGame.IsOver {
		rl.DrawText("Game is over", 190, 200, 20, rl.Red)
	} else {
		offset := int32(0)
		for i := 0; i < MaxHand; i++ {
			if player1.Hand[i] == nil {
				break
			} else {
				if player1.Hand[i] != nil && player1.Hand[i].ToShow == true {
					rl.DrawTexture(player1.Hand[i].Texture, int32(player1.Hand[i].X)+offset, int32(player1.Hand[i].Y), rl.White)
					offset = offset + 30
				}
			}
		}
		rl.DrawText("Score: "+strconv.Itoa(player1.Score), 10, WindowSize-50, 20, rl.RayWhite)
		offset = 0
		isFirst := true
		for i := 0; i < MaxHand; i++ {
			if dealer.Hand[i] == nil {
				break
			} else {
				if dealer.Hand[i] != nil && dealer.Hand[i].ToShow == true {
					if isFirst {
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
	}
	rl.EndDrawing()
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
		getInput(player1, cardDeck, yourGame)
		update(cardDeck, player1, dealer, yourGame)
		draw(background, backOfCard, player1, dealer, yourGame)
		frames++
	}
}
