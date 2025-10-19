package app

import (
	"github.com/LiddleChild/tmux-sessionizer/internal/components/superlist"
	"github.com/LiddleChild/tmux-sessionizer/internal/config"
	"github.com/charmbracelet/lipgloss"
)

var _ superlist.Item = (*entryItem)(nil)

type entryItem config.WorkspaceEntry

func (i entryItem) Label() string {
	return string(i.Path)
}

func (i entryItem) Style(style lipgloss.Style) lipgloss.Style {
	return style
}
