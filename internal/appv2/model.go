package appv2

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/components/superlist"
	"github.com/LiddleChild/tmux-sessionpane/internal/log"
	"github.com/LiddleChild/tmux-sessionpane/internal/tmux"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var _ tea.Model = (*Model)(nil)

type Model struct {
	superlist superlist.Model
}

func New() (Model, error) {
	superlist := superlist.
		New([]superlist.ItemGroup{}).
		SetKeyMap(superlist.KeyMap{
			CursorUp:   keyMap.Up,
			CursorDown: keyMap.Down,
			Submit:     focusedKeyMap.Submit,
			Cancel:     focusedKeyMap.Cancel,
		})

	return Model{
		superlist: superlist,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return tea.Sequence(ListTmuxSessionCmd, SelectAttachedSessionCmd)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Dump(log.LogLevelDebug, msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.superlist, cmd = m.superlist.Update(msg)

		return m, cmd

	case tea.KeyMsg:
		switch {
		case !m.superlist.Focused() && key.Matches(msg, keyMap.Quit):
			return m, tea.Quit

		case !m.superlist.Focused() && key.Matches(msg, keyMap.Select):
			item := m.superlist.GetSelectedItem()

			switch item := item.(type) {
			case *sessionItem:
				execProcessCmd := tea.ExecProcess(tmux.AttachSessionCommand(item.Value()), func(err error) tea.Msg {
					return tea.Sequence(ErrCmd(err), tea.Quit)
				})

				return m, tea.Sequence(
					execProcessCmd,
					tea.Quit,
				)
			}

		case !m.superlist.Focused() && key.Matches(msg, keyMap.Delete):
			item := m.superlist.GetSelectedItem()

			if session, ok := item.(*sessionItem); ok {
				if session.IsAttached {
					return m, nil
				}

				if m.superlist.GetCursor() == len(m.superlist.GetItems())-1 {
					m.superlist.CursorUp()
				}

				if err := tmux.DeleteSession(session.Value()); err != nil {
					return m, ErrCmd(err)
				}

				return m, ListTmuxSessionCmd
			}

		case !m.superlist.Focused() && key.Matches(msg, keyMap.Rename):
			return m, m.superlist.Focus()
		}

	case ListTmuxSessionMsg:
		sessions, err := tmux.ListSession()
		if err != nil {
			return m, tea.Sequence(ErrCmd(err), tea.Quit)
		}

		sessionItems := make([]superlist.Item, 0, len(sessions))
		for _, session := range sessions {
			sessionItems = append(sessionItems, &sessionItem{session})
		}

		entryItems := []superlist.Item{
			entryItem("~/.config/"),
			entryItem("~/dotfiles/"),
			entryItem("~/Documents/Projects/"),
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
		var attached int

	group:
		for _, group := range m.superlist.GetItems() {
			for _, item := range group.Items {
				log.Dump(log.LogLevelDebug, item)
				if item, ok := item.(*sessionItem); ok && item.IsAttached {
					break group
				}

				attached += 1
			}
		}

		m.superlist.SetCursor(attached)
		return m, nil

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
	return m.superlist.View()
}
