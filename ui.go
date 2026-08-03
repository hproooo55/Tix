package main

import (
	// "charm.land/bubbles/v2/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	headStyle = lipgloss.NewStyle()
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

	}

	v := s
	return v
}
