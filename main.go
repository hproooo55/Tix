package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/wish/bubbletea"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	wish "github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
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
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) { return teaHandler(s, &l) }),
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

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	m := model{
		term:      pty.Term,
		width:     pty.Window.Width,
		height:    pty.Window.Height,
		player:    player,
		nameInput: ti,
		players:   l.players,
		txtStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		quitStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		bg:        "light",
	}
	return m, []tea.ProgramOption{}
}

const (
	titleView uint = iota
	LobbyView
	GameView
)

// Just a generic tea.Model to demo terminal information of ssh.
type model struct {
	term       string
	width      int
	height     int
	player     *Player
	players    []*Player
	bg         string
	view       uint
	txtStyle   lipgloss.Style
	quitStyle  lipgloss.Style
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
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.player.displayName = m.nameInput.Value()
		}
		switch m.view {
		case titleView:
			if m.player.displayName == "" {
				var cmd tea.Cmd
				m.nameInput, cmd = m.nameInput.Update(msg)
				return m, cmd
			}
		}

	}
	return m, nil
}
