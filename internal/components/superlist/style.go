package superlist

import "github.com/charmbracelet/lipgloss"

var (
	groupNameStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("8"))

	hoverItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("5")).
			Background(lipgloss.Color("8")).
			Bold(true)
)
