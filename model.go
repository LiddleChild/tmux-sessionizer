package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/internal/listinput"
	"github.com/LiddleChild/tmux-sessionpane/internal/tmux"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
)

var _ listinput.Item = (*sessionItem)(nil)

type sessionItem tmux.Session

func (i sessionItem) Label() string {
	if i.IsAttached {
		return i.Name + " (attached)"
	}

	return i.Name
}

func (i sessionItem) Value() string {
	return i.Name
}

func (i sessionItem) FilterValue() string {
	return i.Name
}

var _ tea.Model = (*model)(nil)

type model struct {
	dump io.Writer
	err  error

	keys keyMap
	help help.Model

	list listinput.Model
}

func NewModel(dump io.Writer) (*model, error) {
	// sessions, err := tmux.ListSession()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// items := []listinput.Item{}
	// for _, session := range sessions {
	// 	items = append(items, sessionItem(session))
	// }
	//
	l := listinput.New([]listinput.Item{}, 0, 0)

	l.SetKeyMap(list.KeyMap{
		CursorUp:   keymap.Up,
		CursorDown: keymap.Down,
	})

	return &model{
		dump: dump,
		err:  nil,
		keys: keymap,
		help: help.New(),
		list: l,
	}, nil
}

func (m model) Init() tea.Cmd {
	return ListTmuxSessionCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
	}

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

	case listinput.InputSubmitedMsg:
		item := m.list.Items()[msg.Index]

		if err := tmux.RenameSession(item.Value(), msg.Value); err != nil {
			return m, QuitWithErr(err)
		}

		return m, ListTmuxSessionCmd

	case ListTmuxSessionMsg:
		sessions, err := tmux.ListSession()
		if err != nil {
			return m, QuitWithErr(err)
		}

		items := []listinput.Item{}
		for _, session := range sessions {
			items = append(items, sessionItem(session))
		}

		return m, m.list.SetItems(items)
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
