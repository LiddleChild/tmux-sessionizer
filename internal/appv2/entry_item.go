package appv2

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/config"
	"github.com/charmbracelet/lipgloss"
)

type entryItem config.WorkspaceEntry

func (i entryItem) Label() string {
	return string(i.Path)
}

func (i entryItem) Style(style lipgloss.Style) lipgloss.Style {
	return style
}
