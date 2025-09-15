package main

import (
	"fmt"
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/tmux"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ list.Item = (*item)(nil)

type item tmux.Session

func (i item) Title() string {
	if i.IsAttached {
		return i.Name + " (attached)"
	}

	return i.Name
}

func (i item) FilterValue() string {
	return i.Name
}

var _ tea.Model = (*model)(nil)

type model struct {
	err error

	keys keyMap
	help help.Model

	list list.Model
}

func NewModel() (*model, error) {
	sessions, err := tmux.ListSession()
	if err != nil {
		return nil, err
	}

	items := []list.Item{}
	for _, session := range sessions {
		items = append(items, item(session))
	}

	l := list.New(items, itemDelegate{}, 0, len(sessions))
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	l.KeyMap = list.KeyMap{
		CursorUp:   keymap.Up,
		CursorDown: keymap.Down,
	}

	return &model{
		err:  nil,
		keys: keymap,
		help: help.New(),
		list: l,
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, keymap.Select):
			session := m.list.SelectedItem().(item)

			execProcessCmd := tea.ExecProcess(tmux.AttachSessionCommand(session.Name), func(err error) tea.Msg {
				return QuitWithErr(err)
			})

			return m, tea.Sequence(
				execProcessCmd,
				tea.Quit,
			)
		}

	case QuitWithErrMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s %s", AppName, Version))
	builder.WriteByte('\n')
	builder.WriteString(m.help.View(keymap))
	builder.WriteByte('\n')
	builder.WriteByte('\n')
	builder.WriteString(m.list.View())

	return builder.String()
}
