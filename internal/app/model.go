package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LiddleChild/tmux-sessionizer/internal/components/superlist"
	"github.com/LiddleChild/tmux-sessionizer/internal/config"
	"github.com/LiddleChild/tmux-sessionizer/internal/log"
	"github.com/LiddleChild/tmux-sessionizer/internal/tmux"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = (*Model)(nil)

type Model struct {
	superlist superlist.Model
	help      help.Model
}

func New() (Model, error) {
	superlist := superlist.
		New([]superlist.ItemGroup{}).
		SetKeyMap(superlist.KeyMap{
			CursorUp:   keyMap.Up,
			CursorDown: keyMap.Down,
			FocusItem:  keyMap.Rename,
			Submit:     focusedKeyMap.Submit,
			Cancel:     focusedKeyMap.Cancel,
		})

	return Model{
		superlist: superlist,
		help:      help.New(),
	}, nil
}

func (m Model) renderTopBar() string {
	var help string
	if m.superlist.FocusedComponent() == superlist.FocusedComponentItem {
		help = m.help.FullHelpView(focusedKeyMap.FullHelp())
	} else {
		help = m.help.FullHelpView(keyMap.FullHelp())
	}

	return lipgloss.NewStyle().
		Height(4).
		Render(lipgloss.JoinVertical(lipgloss.Top,
			fmt.Sprintf("%s %s", config.AppName, config.AppVersion),
			help,
		))
}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(ListTmuxSessionCmd, SelectAttachedSessionCmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Dump(log.LogLevelDebug, msg)

	switch msg := msg.(type) {
	case ErrMsg:
		log.Println(log.LogLevelError, strings.TrimSpace(msg.Error()))

	case tea.WindowSizeMsg:
		helpHeight := lipgloss.Height(m.renderTopBar())

		var cmd tea.Cmd
		m.superlist, cmd = m.superlist.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - helpHeight,
		})

		return m, cmd

	case tea.KeyMsg:
		switch {
		case m.superlist.FocusedComponent() == superlist.FocusedComponentFilter && key.Matches(msg, keyMap.Quit):
			return m, tea.Quit

		case m.superlist.FocusedComponent() == superlist.FocusedComponentFilter && key.Matches(msg, keyMap.Select):
			switch item := m.superlist.GetSelectedItem().(type) {
			case *sessionItem:
				execProcessCmd := tea.ExecProcess(tmux.AttachSessionCommand(item.Name), func(err error) tea.Msg {
					return tea.Sequence(ErrCmd(err), tea.Quit)
				})

				return m, tea.Sequence(
					execProcessCmd,
					tea.Quit,
				)

			case *entryItem:
				if !tmux.HasSession(item.Name) {
					if err := tmux.NewDetachedSession(item.Name, item.Path); err != nil {
						return m, tea.Sequence(ErrCmd(err), tea.Quit)
					}
				}

				execProcessCmd := tea.ExecProcess(tmux.AttachSessionCommand(item.Name), func(err error) tea.Msg {
					return tea.Sequence(ErrCmd(err), tea.Quit)
				})

				return m, tea.Sequence(
					execProcessCmd,
					tea.Quit,
				)

			default:
				return m, nil
			}

		case m.superlist.FocusedComponent() == superlist.FocusedComponentFilter && key.Matches(msg, keyMap.Delete):
			item := m.superlist.GetSelectedItem()
			if session, ok := item.(*sessionItem); ok {
				if session.IsAttached {
					return m, nil
				}

				m.superlist.CursorUp()

				if err := tmux.DeleteSession(session.Value()); err != nil {
					return m, ErrCmd(err)
				}

				return m, ListTmuxSessionCmd
			}
		}

	case ListTmuxSessionMsg:
		sessions, err := tmux.ListSessions()
		if errors.Is(err, tmux.NoServerRunningErr) {
			sessions = []tmux.Session{}
		} else if err != nil {
			return m, tea.Sequence(ErrCmd(err), tea.Quit)
		}

		sessionItems := make([]superlist.Item, 0, len(sessions))
		for _, session := range sessions {
			si := sessionItem(session)
			sessionItems = append(sessionItems, &si)
		}

		entryItems := make([]superlist.Item, 0, len(config.WorkspaceEntries))
		for _, entry := range config.WorkspaceEntries {
			ei := entryItem(entry)
			entryItems = append(entryItems, &ei)
		}

		items := []superlist.ItemGroup{
			{
				Name:  "Sessions",
				Items: sessionItems,
			},
			{
				Name:  "Entries",
				Items: entryItems,
			},
		}

		m.superlist.SetItems(items)

		return m, nil

	case SelectAttachedSessionMsg:
		for i, item := range m.superlist.GetItemIter() {
			if item, ok := item.(*sessionItem); ok && item.IsAttached {
				m.superlist.SetCursor(i)
				return m, nil
			}
		}

		m.superlist.SetCursor(0)
		return m, nil

	case superlist.FilterMsg:
		var cmd tea.Cmd
		m.superlist, cmd = m.superlist.Update(msg)
		return m, cmd

	case superlist.SubmitMsg:
		if err := tmux.RenameSession(msg.OldValue, msg.NewValue); err != nil {
			return m, ErrCmd(err)
		}

		return m, ListTmuxSessionCmd
	}

	var cmd tea.Cmd
	m.superlist, cmd = m.superlist.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.renderTopBar(),
		m.superlist.View(),
	)
}
