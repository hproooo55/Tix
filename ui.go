package main

import (
	"charm.land/lipgloss/v2"
)

var (
	headStyle = lipgloss.NewStyle().
		MarginLeft(50).
		MarginRight(50).
		PaddingLeft(1).
		PaddingRight(1)
)

func (m model) View() string {
	colors := lipgloss.Blend1D(15, lipgloss.Color("#b5e3ff"), lipgloss.Color("#0167ff"), lipgloss.Color("#b5e3ff"))

	m.headStyle = headStyle.Foreground(colors[m.colorIndex])
	v := m.headStyle.Render(`
__/\\\\\\\\\\\\\\\__/\\\\\\\\\\\__/\\\_______/\\\_        
 _\///////\\\/////__\/////\\\///__\///\\\___/\\\/__       
  _______\/\\\___________\/\\\_______\///\\\\\\/____      
   _______\/\\\___________\/\\\_________\//\\\\______     
    _______\/\\\___________\/\\\__________\/\\\\______    
     _______\/\\\___________\/\\\__________/\\\\\\_____   
      _______\/\\\___________\/\\\________/\\\////\\\___  
       _______\/\\\________/\\\\\\\\\\\__/\\\/___\///\\\_ 
        _______\///________\///////////__\///_______\///__
		`)
	return v
}
