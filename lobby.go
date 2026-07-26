package main

import (
	"sync"
)

type Mark byte

const (
	Empty Mark = 0
	Xmark Mark = 'X'
	Omark Mark = 'O'
)

func (m Mark) String() string {
	if m == Empty {
		return ""
	}
	return string(rune(m))
}

func (m Mark) opposite() Mark {
	if m == Omark {
		return Xmark
	}

	return Omark
}

type Status string

const (
	statusPlaying  Status = "waiting"
	statusEnded    Status = "ended"
	statusProgress Status = "in progress"
)

var winLines = [8][3]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
	{0, 4, 8}, {2, 4, 6},
}

// define lobby structure
type Lobby struct {
	mu   sync.Mutex
	turn *Player
}

type Player struct {
	game        *Game
	mark        Mark
	displayName string
}

type Game struct {
	ID    string
	seq   int64
	Board [9]Mark
	turn  Mark
	X, O  *Player
}
