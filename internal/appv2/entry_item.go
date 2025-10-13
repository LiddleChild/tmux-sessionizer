package appv2

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/config"
	"github.com/charmbracelet/lipgloss"
)

type entryItem struct {
	config.WorkspaceEntry
}

func (i entryItem) Name() string {
	return string(i.WorkspaceEntry.Path)
}

func (i entryItem) Style(style lipgloss.Style) lipgloss.Style {
	return style
}
