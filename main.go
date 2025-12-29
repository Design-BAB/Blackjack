//Author: Design-BAB
//Date: 12/25/2025
//Description: Time to play Blackjack! The goal is to reach 452 lines of code

package main

import (
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
)

type GameState struct {
	JustStarted   bool
	IsOver        bool
	YourHandValue int
	HerHandValue  int
}

func newGame() *GameState {
	return &GameState{JustStarted: true}
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
			cardDeck[c] = newCard(cardTextures[c], Center, WindowSize-200, 10)
		case 14:
			cardDeck[c] = newCard(cardTextures[c], Center, WindowSize-200, 11)
		default:
			cardDeck[c] = newCard(cardTextures[c], Center, WindowSize-200, i)
		}
		//ha ha ha no, this is not a programming pun. c refers to card count therefore c++
		c++
	}
	return cardTextures, cardDeck, c
}

func getInput(yourGame *GameState) {
	if yourGame.IsOver == false {

	}
}

func update(cardDeck [TotalDeck]*Card, yourHand, herHand [MaxHand]*CardInHand, yourGame *GameState) ([MaxHand]*CardInHand, [MaxHand]*CardInHand) {
	if yourGame.JustStarted {
		rand.Shuffle(TotalDeck, func(i, j int) {
			cardDeck[i], cardDeck[j] = cardDeck[j], cardDeck[i]
		})
		//giving first card
		yourHand[0] = AddCardToHand(cardDeck[0])
		yourHand[0].ToShow = true
		cardDeck[0].IsDiscarded = true
		//giving second card
		yourHand[1] = AddCardToHand(cardDeck[1])
		yourHand[1].ToShow = true
		cardDeck[1].IsDiscarded = true
		//giving first card to dealer
		herHand[0] = AddCardToHand(cardDeck[2])
		herHand[0].ToShow = true
		herHand[0].X = WindowSize - 5 - herHand[0].Width
		herHand[0].Y = 50
		cardDeck[2].IsDiscarded = true
		//giving second card to dealer
		herHand[1] = AddCardToHand(cardDeck[3])
		herHand[1].ToShow = true
		herHand[1].X = WindowSize - 5 - herHand[1].Width
		herHand[1].Y = 50
		cardDeck[3].IsDiscarded = true
		yourGame.JustStarted = false
	}
	return yourHand, herHand
}

func draw(background, backOfCard rl.Texture2D, yourHand, herHand [MaxHand]*CardInHand, yourGame *GameState) {
	rl.BeginDrawing()

	rl.ClearBackground(rl.RayWhite)
	//rl.DrawText("Congrats! You created your first window!", 190, 200, 20, rl.LightGray)
	rl.DrawTexture(background, 0, 0, rl.White)
	if yourGame.IsOver {
		rl.DrawText("Game is over", 190, 200, 20, rl.Red)
	}
	var offset int32
	offset = 0
	isFirst := true
	for _, card := range yourHand {
		if card != nil && card.ToShow == true {
			if isFirst {
				yourGame.YourHandValue = card.Value
			} else {
				yourGame.YourHandValue = yourGame.YourHandValue + card.Value
			}
			rl.DrawTexture(card.Texture, int32(card.X)+offset, int32(card.Y), rl.White)
			offset = offset + 30
		}
	}
	offset = 0
	isFirst = true
	for _, card := range herHand {
		if card != nil && card.ToShow == true {
			if isFirst {
				rl.DrawTexture(backOfCard, int32(card.X), int32(card.Y), rl.White)
				yourGame.HerHandValue = card.Value
				offset = offset + 40
				isFirst = false
			} else {
				rl.DrawTexture(card.Texture, int32(card.X), int32(card.Y)+offset, rl.White)
				offset = offset + 40
				yourGame.HerHandValue = yourGame.HerHandValue + card.Value
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
	cardTextures, cardDeck := importCards()
	var yourHand [MaxHand]*CardInHand
	var herHand [MaxHand]*CardInHand
	for _, texture := range cardTextures {
		defer rl.UnloadTexture(texture)
	}
	frames := 0
	for !rl.WindowShouldClose() && frames < MaxFrames {
		getInput(yourGame)
		yourHand, herHand = update(cardDeck, yourHand, herHand, yourGame)
		draw(background, backOfCard, yourHand, herHand, yourGame)
		frames++
	}
}
