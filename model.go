package main

import (
	"fmt"
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/internal/listinput"
	"github.com/LiddleChild/tmux-sessionpane/tmux"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var _ listinput.Item = (*sessionItem)(nil)

type sessionItem struct {
	session tmux.Session
}

func (i *sessionItem) Name() string {
	if i.session.IsAttached {
		return i.session.Name + " (attached)"
	}

	return i.session.Name
}

func (i *sessionItem) Value() string {
	return i.session.Name
}

func (i *sessionItem) OnValueChange(value string) tea.Cmd {
	i.session.Name = value
	return func() tea.Msg {
		if err := tmux.RenameSession(i.session.Name, value); err != nil {
			return QuitWithErr(err)
		}

		return nil
	}
}

func (i *sessionItem) FilterValue() string {
	return i.session.Name
}

var _ tea.Model = (*model)(nil)

type model struct {
	err error

	keys keyMap
	help help.Model

	list listinput.Model
}

func NewModel() (*model, error) {
	sessions, err := tmux.ListSession()
	if err != nil {
		return nil, err
	}

	items := []listinput.Item{}
	for _, session := range sessions {
		items = append(items, &sessionItem{session})
	}

	l := listinput.New(items, 0, len(sessions))
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
		case !m.list.IsFocused() && key.Matches(msg, keymap.Quit):
			return m, tea.Quit

		case !m.list.IsFocused() && key.Matches(msg, keymap.Select):
			session := m.list.SelectedItem().(listinput.Item)

			execProcessCmd := tea.ExecProcess(tmux.AttachSessionCommand(session.Name()), func(err error) tea.Msg {
				return QuitWithErr(err)
			})

			return m, tea.Sequence(
				execProcessCmd,
				tea.Quit,
			)

		case !m.list.IsFocused() && key.Matches(msg, keymap.Rename):
			return m, m.list.FocusSelectedItem()
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
