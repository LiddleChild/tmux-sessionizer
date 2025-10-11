package appv2

import (
	"fmt"

	"github.com/LiddleChild/tmux-sessionpane/internal/components/superlist"
	"github.com/LiddleChild/tmux-sessionpane/internal/tmux"
	tea "github.com/charmbracelet/bubbletea"
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

func (i *sessionItem) SetValue(value string) tea.Cmd {
	if err := tmux.RenameSession(i.Session.Name, value); err != nil {
		return ErrCmd(err)
	}

	i.Session.Name = value
	return ListTmuxSessionCmd
}
