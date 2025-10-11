package appv2

import (
	"fmt"

	"github.com/LiddleChild/tmux-sessionpane/internal/components/superlist"
	"github.com/LiddleChild/tmux-sessionpane/internal/tmux"
	"github.com/charmbracelet/lipgloss"
)

var _ superlist.InputItem = (*sessionItem)(nil)

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

func (i sessionItem) Value() string {
	return i.Session.Name
}

func (i *sessionItem) SetValue(value string) {
	i.Session.Name = value
}
