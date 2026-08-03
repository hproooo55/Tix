package main

import (
	"errors"
	"log"
	"slices"
	"sync"

	"github.com/google/uuid"
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

// Shared app state between multiple players
type Lobby struct {
	mu      sync.Mutex
	games   []*Game
	players []*Player
}

type Player struct {
	game        *Game
	mark        Mark
	displayName string
	remoteAddr  string
}

type Game struct {
	ID     string
	locked bool
	// seq   int64
	Board [9]Mark
	turn  Mark
	X, O  *Player
}

func (l *Lobby) getPlayerByAddr(addr string) *Player {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := range len(l.players) {
		if l.players[i].remoteAddr == addr {
			return l.players[i]
		}
	}
	return nil
}

func (l *Lobby) createGame(p *Player) (*Game, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, e := range l.games {
		if e.X == p || e.O == p {
			// delete game directly inside createGame to avoid mutex handling.
			for i := range l.games {
				if l.games[i] == e {
					// Do not allow a user to create more than 1 game.
					l.games = slices.Delete(l.games, i, i+1)
					log.Println(i)
				}
			}
		}
	}
	id := uuid.New()
	g := Game{
		ID:    id.String(),
		Board: [9]Mark{},
		turn:  Xmark,
		X:     p,
	}
	l.games = append(l.games, &g)
	p.mark = Xmark
	return &g, nil
}

// func (l *Lobby) updateGameState(p *Player) (*Game, error) {

// }

func (l *Lobby) deleteGame(g *Game) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i := range l.games {
		if l.games[i] == g {
			l.games = slices.Delete(l.games, i, i+1)
			log.Println(i)
			log.Println("Delted")
			return nil
		}
	}

	return errors.New("GAME NOT FOUND")
}
