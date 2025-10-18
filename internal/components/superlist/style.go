package superlist

import "github.com/charmbracelet/lipgloss"

var (
	groupNameStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("8"))

	hoveredItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("6")).
				Background(lipgloss.Color("8")).
				Bold(true)

	filterPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("6")).
				Bold(true)
)
