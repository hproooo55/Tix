package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	headStyle         = lipgloss.NewStyle()
	selectedCellStyle = lipgloss.NewStyle()
	CellStyle         = lipgloss.NewStyle()
	hor               = lipgloss.NewStyle()
	rows              = 3
	col               = 3
)

func (m model) View() string {
	m.headStyle = headStyle.Align(lipgloss.Center).Foreground(lipgloss.Color("#46B1C9")).Width(m.width)
	// menuStyles := []lipgloss.Style{
	// 	lipgloss.NewStyle().Foreground(lipgloss.Color("#4A6FA5")),
	// }
	menuItems := []string{
		"Create Game",
		"List Games",
		"Join Game By Code",
	}

	var ms string
	var mo string
	prefix := lipgloss.NewStyle().Foreground(lipgloss.Color("#A0DDFF")).Render("> ")
	if m.player.displayName != "" {
		for i, mi := range menuItems {
			if i == m.menuChoice {
				ms = ms + lipgloss.NewStyle().Align(lipgloss.Left).Render(prefix+mi) + "\n"
			} else {
				ms = ms + lipgloss.NewStyle().Align(lipgloss.Left).Foreground(lipgloss.Color("#4A6FA5")).Render(mi) + "\n"
			}
			mo = mo + lipgloss.PlaceHorizontal(m.width, lipgloss.Center, ms)
		}
	}

	ascii := `
__/\\\\\\\\\\\\\\\__/\\\\\\\\\\\__/\\\_______/\\\_        
 _\///////\\\/////__\/////\\\///__\///\\\___/\\\/__       
  _______\/\\\___________\/\\\_______\///\\\\\\/____      
   _______\/\\\___________\/\\\_________\//\\\\______     
    _______\/\\\___________\/\\\__________\/\\\\______    
     _______\/\\\___________\/\\\__________/\\\\\\_____   
      _______\/\\\___________\/\\\________/\\\////\\\___  
       _______\/\\\________/\\\\\\\\\\\__/\\\/___\///\\\_ 
        _______\///________\///////////__\///_______\///__ 
	`
	var ps string
	for _, p := range m.lobby.players {
		ps = ps + p.displayName + "\n"
	}
	var wm string
	var s string
	switch m.view {
	case titleView:
		if m.player.displayName != "" && m.view == titleView {
			wm = "WELCOME! " + m.player.displayName
		} else {
			wm = m.nameInput.View()
		}
		s = ascii + "\n\n" + wm + "\n\n" + ms
		s = m.headStyle.Render(s)
	case LobbyView:
		s = m.table.View()
		m.headStyle = headStyle.Foreground(lipgloss.Color("#46B1C9")).Width(m.table.Width()).Border(lipgloss.NormalBorder())

		tbs := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.headStyle.Render(s))
		s = tbs
	case GameView:
		var str string = ""
		var vert string
		for i := range col {
			vert = ""
			for j := 0; j < 3; j++ {
				var c string
				idx := i*3 + j
				if m.player.selectedCell == idx {
					c = CellStyle.Height(5).Width(10).Border(lipgloss.BlockBorder(), true, true, true, true).Background(lipgloss.Color("#46B1C9")).Align(lipgloss.Center).Padding(1).Render(string(m.player.game.Board[i*3+j]))

				} else {
					c = CellStyle.Height(5).Width(10).Border(lipgloss.BlockBorder(), true, true, true, true).Align(lipgloss.Center).Padding(1).Render(string(m.player.game.Board[i*3+j]))

				}
				vert = lipgloss.JoinHorizontal(lipgloss.Center, vert, c)
			}
			str = lipgloss.JoinVertical(lipgloss.Center, str, vert)
		}

		// handle O not existing, the same thing is impossible for X (the host)
		var ODisplayName string
		if m.player.game.O.displayName == "" {
			ODisplayName = "Waiting for opponent..."
		} else {
			ODisplayName = m.player.displayName
		}
		s = lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.player.game.X.displayName+"\t vs \t"+ODisplayName+"\n\n"+str)
	}

	v := s
	return v
}
