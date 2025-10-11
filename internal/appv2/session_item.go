package appv2

import (
	"fmt"

	"github.com/LiddleChild/tmux-sessionpane/internal/tmux"
	"github.com/charmbracelet/lipgloss"
)

type sessionItem struct {
	tmux.Session
}

func (i sessionItem) Name() string {
	name := i.Session.Name
	if i.IsAttached {
		name = fmt.Sprintf("%s (attached)", name)
	}

	return name
}

func (i sessionItem) Style(style lipgloss.Style) lipgloss.Style {
	if i.IsAttached {
		return style.Bold(true)
	}

	return style
}
