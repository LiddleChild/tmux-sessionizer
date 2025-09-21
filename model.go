package main

import (
	"fmt"
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/internal/config"
	"github.com/LiddleChild/tmux-sessionpane/internal/listinput"
	"github.com/LiddleChild/tmux-sessionpane/internal/log"
	"github.com/LiddleChild/tmux-sessionpane/internal/tmux"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ listinput.Item = (*sessionItem)(nil)

type sessionItem tmux.Session

func (i sessionItem) Label() string {
	if i.IsAttached {
		return i.Name + " (attached)"
	}

	return i.Name
}

func (i sessionItem) Style(style lipgloss.Style) lipgloss.Style {
	if i.IsAttached {
		return style.Bold(true)
	}

	return style
}

func (i sessionItem) Value() string {
	return i.Name
}

func (i *sessionItem) SetValue(value string) tea.Cmd {
	if err := tmux.RenameSession(i.Name, value); err != nil {
		return ErrCmd(err)
	}

	i.Name = value
	return ListTmuxSessionCmd
}

func (i sessionItem) FilterValue() string {
	return i.Name
}

var _ tea.Model = (*model)(nil)

type model struct {
	keys keyMap
	help help.Model

	list listinput.Model
}

func NewModel() (*model, error) {
	l := listinput.New([]listinput.Item{})

	l.SetKeyMap(listinput.KeyMap{
		CursorUp:   keymap.Up,
		CursorDown: keymap.Down,
		Submit:     focusedKeymap.Submit,
		Cancel:     focusedKeymap.Cancel,
	})

	return &model{
		keys: keymap,
		help: help.New(),
		list: l,
	}, nil
}

func (m model) Init() tea.Cmd {
	return ListTmuxSessionCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Dump(log.LogLevelDebug, msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch {
		case !m.list.IsFocused() && key.Matches(msg, keymap.Quit):
			return m, tea.Quit

		case !m.list.IsFocused() && key.Matches(msg, keymap.Select):
			session := m.list.SelectedItem()

			execProcessCmd := tea.ExecProcess(tmux.AttachSessionCommand(session.Value()), func(err error) tea.Msg {
				return tea.Sequence(ErrCmd(err), tea.Quit)
			})

			return m, tea.Sequence(
				execProcessCmd,
				tea.Quit,
			)

		case !m.list.IsFocused() && key.Matches(msg, keymap.Rename):
			return m, m.list.FocusSelectedItem()

		case !m.list.IsFocused() && key.Matches(msg, keymap.Delete):
			session := m.list.SelectedItem().(*sessionItem)

			if session.IsAttached {
				return m, nil
			}

			if m.list.Index() == len(m.list.Items())-1 {
				m.list.CursorUp()
			}

			if err := tmux.DeleteSession(session.Name); err != nil {
				return m, ErrCmd(err)
			}

			return m, ListTmuxSessionCmd
		}

	case ErrMsg:
		log.Printlnf(log.LogLevelError, msg.Error())
		return m, nil

	case ListTmuxSessionMsg:
		sessions, err := tmux.ListSession()
		if err != nil {
			return m, tea.Sequence(ErrCmd(err), tea.Quit)
		}

		items := make([]listinput.Item, 0, len(sessions))
		for _, session := range sessions {
			sessionItem := sessionItem(session)
			items = append(items, &sessionItem)
		}

		return m, m.list.SetItems(items)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s %s", config.AppName, config.Version))
	builder.WriteByte('\n')

	if m.list.IsFocused() {
		builder.WriteString(m.help.FullHelpView(focusedKeymap.FullHelp()))
		builder.WriteByte('\n')
	} else {
		builder.WriteString(m.help.FullHelpView(keymap.FullHelp()))
	}

	builder.WriteByte('\n')
	builder.WriteByte('\n')
	builder.WriteString(m.list.View())

	return builder.String()
}
