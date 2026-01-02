🃏 Blackjack (Raylib + Go)

A desktop Blackjack game built in Go using raylib-go, featuring full player vs dealer gameplay, card animations, scoring logic, and round-based progression with lives.

This project was created as a hands-on way to deepen understanding of game loops, state management, and real-time input handling in Go.

🎮 Gameplay Features

Classic Blackjack rules

Player actions: Hit and Stay

Dealer logic (hits until 17)

Automatic Ace value adjustment (11 → 1)

Round-based gameplay with lives system

Card shuffling and deck reset

Visual card rendering with hidden dealer card

Timed dealer actions for realism

🛠️ Built With

Go

raylib-go

2D sprite-based rendering

Real-time game loop with frame control

🎯 Controls
Key	Action
X	Hit
Z	Stay
A	Restart game (after Game Over)
📁 Project Structure (High Level)

GameState
Manages turns, rounds, lives, and overall game flow

Player / Dealer
Handles hands, scores, and role-based behavior

Card / Deck Logic
Card values, discard tracking, shuffling, and dealing

Rendering
Separate drawing functions for cards, UI, and game messages



📸 Screenshots

(Add gameplay screenshots or a short GIF here — highly recommended for portfo
