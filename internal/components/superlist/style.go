package superlist

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/colors"
	"github.com/charmbracelet/lipgloss"
)

var (
	groupNameStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(colors.BrightBlack)

	hoveredItemStyle = lipgloss.NewStyle().
				Foreground(colors.Cyan).
				Background(colors.BrightBlack).
				Bold(true)

	filterPromptStyle = lipgloss.NewStyle().
				Foreground(colors.BrightBlack).
				Bold(true)
)
