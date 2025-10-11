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
	sessions, err := tmux.ListSession()
	if err != nil {
		return Model{}, err
	}

	sessionItems := make([]superlist.Item, 0, len(sessions))
	for _, s := range sessions {
		sessionItems = append(sessionItems, sessionItem{s})
	}

	entryItems := []superlist.Item{
		entryItem("~/.config/"),
		entryItem("~/dotfiles/"),
		entryItem("~/Documents/Projects/"),
	}

	groups := []superlist.ItemGroup{
		{
			Name:  "Sessions",
			Items: sessionItems,
		},
		{
			Name:  "Entries",
			Items: entryItems,
		},
	}

	superlist := superlist.New(groups).SetKeyMap(superlist.KeyMap{
		CursorUp:   keyMap.Up,
		CursorDown: keyMap.Down,
	})

	return Model{
		superlist: superlist,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Dump(log.LogLevelDebug, msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.superlist, cmd = m.superlist.Update(msg)

		return m, cmd

	case tea.KeyMsg:
		if key.Matches(msg, keyMap.Quit) {
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.superlist, cmd = m.superlist.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return m.superlist.View()
}
