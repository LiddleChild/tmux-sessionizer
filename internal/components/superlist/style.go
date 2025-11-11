package superlist

import (
	"github.com/LiddleChild/tmux-sessionizer/internal/colors"
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

	noResultStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(colors.BrightBlack).
			AlignHorizontal(lipgloss.Center)
)
