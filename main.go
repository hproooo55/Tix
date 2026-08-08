package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	wishbubbletea "github.com/charmbracelet/wish/bubbletea"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	wish "github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "127.0.0.1"
	port = "22"
)

func main() {
	l := Lobby{
		games:   []*Game{},
		players: []*Player{},
	}
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			wishbubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) { return teaHandler(s, &l) }),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatal("Could not start ssh server", "error", err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping ssh server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}

}

type Disconnected struct{}

// func listenDisconnect(ctx ) tea.Cmd {
// 	return func() tea.Msg{
// 		c
// 		return Disconnected
// 	}
// }

func teaHandler(s ssh.Session, l *Lobby) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	addr := strings.Split(s.RemoteAddr().String(), ":")[0]

	player := l.getPlayerByAddr(addr)
	if player == nil {
		player = &Player{
			remoteAddr: addr,
		}
		l.players = append(l.players, player)
	}

	log.Info(player.remoteAddr)
	log.Info(addr)

	// directly proportional to length of games + 10 extra room
	h := pty.Window.Height * (len(l.games)*5 + 10) / 100
	log.Info(h)

	width := pty.Window.Width * 30 / 100
	height := h

	tb := table.New(
		table.WithFocused(true),
		table.WithHeight(height),
		table.WithWidth(width),
	)

	tb.SetColumns([]table.Column{
		{
			Title: "Index",
			Width: width * 25 / 100,
		},
		{
			Title: "Player",
			Width: width * 50 / 100,
		},
	})
	ts := table.DefaultStyles()

	ts.Header = ts.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("#46B1C9")).
		Bold(true).
		Foreground(lipgloss.Color("#fff"))

	ts.Selected = ts.Selected.Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false).
		Width(width)
	tb.SetStyles(ts)

	ti := textinput.New()
	ti.Placeholder = "Enter your display name here"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	m := model{
		term:       pty.Term,
		width:      pty.Window.Width,
		height:     pty.Window.Height,
		player:     player,
		ctx:        s.Context(),
		lobby:      l,
		nameInput:  ti,
		listChoice: 0,
		menuChoice: 0,
		table:      tb,
		txtStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		quitStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}

	go func() {
		<-s.Context().Done()
		// delete game if any
		if m.player.game != nil {
			l.deleteGame(m.player.game)
		}
		log.Info("Client disconnected")

	}()

	return m, []tea.ProgramOption{}
}

const (
	titleView uint = iota
	LobbyView
	// we dont need createview rn bcz there isnt any input and not necessary to show a view
	// CreateView
	JoinView
	GameView
)

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term       string
	width      int
	height     int
	player     *Player
	menuChoice int
	lobby      *Lobby
	ctx        context.Context
	view       uint
	listChoice int
	txtStyle   lipgloss.Style
	quitStyle  lipgloss.Style
	table      table.Model
	headStyle  lipgloss.Style
	nameInput  textinput.Model
	colorIndex int
}

func (m model) Init() tea.Cmd {
	return tick()
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
func getRows(l *Lobby) []table.Row {
	tr := []table.Row{}
	for i, g := range l.games {
		tr = append(tr, table.Row{
			strconv.Itoa(i),
			g.X.displayName,
		})
	}
	return tr
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tickMsg:
		m.colorIndex = (m.colorIndex + 1) % 15
		return m, tick()
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "w":
			if m.menuChoice < 0 {
				m.menuChoice = 2
			} else {
				m.menuChoice--
			}

		case "down", "s":
			if m.menuChoice > 2 {
				m.menuChoice = 0
			} else {
				m.menuChoice++
			}

		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.view == titleView && m.player.displayName != "" {
				switch m.menuChoice {
				case 0:
					p, err := m.lobby.createGame(m.player)
					if err == nil {
						m.player = p
						m.view = GameView
					}
				case 1:
					m.view = LobbyView
				case 2:
					m.view = JoinView
				}
			} else if m.view == titleView && m.player.displayName == "" {
				m.player.displayName = m.nameInput.Value()
			} else if m.view == LobbyView {
				var gameToJoin *Game
				for _, game := range m.lobby.games {
					if m.table.SelectedRow()[1] == game.X.displayName {
						gameToJoin = game
					}
				}
				p, err := m.lobby.joinGame(m.player, gameToJoin)
				m.player = p
				if err == nil {
					m.view = GameView
				}
			} else if m.view == GameView {

			}
		}
		switch m.view {
		case titleView:
			if m.player.displayName == "" {
				var cmd tea.Cmd
				m.nameInput, cmd = m.nameInput.Update(msg)
				return m, cmd
			}
		case LobbyView:
			var cmd tea.Cmd
			m.table, cmd = m.table.Update(msg)
			m.table.SetRows(getRows(m.lobby))
			return m, cmd

		case GameView:
			if m.player.game == nil {
				log.Print("SETTING LObbY VIEW")
				m.view = LobbyView
			}
			switch msg.String() {
			case "down":
				if m.player.selectedCell > 6 {
					m.player.selectedCell -= 9
				} else {
					m.player.selectedCell += 3
				}
			case "up":
				if m.player.selectedCell < 0 {
					m.player.selectedCell += 9
				} else {
					m.player.selectedCell -= 3
				}
			case "right":
				if m.player.selectedCell > 8 {
					m.player.selectedCell -= 9
				} else {
					m.player.selectedCell++
				}
			case "left":
				if m.player.selectedCell < 0 {
					m.player.selectedCell += 9
				} else {
					m.player.selectedCell--
				}
			case "enter":
				m.player.game.makeMove(m.player)
			}
		}
	}
	return m, nil
}
