package superlist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	groups []ItemGroup
	cursor int

	width  int
	height int

	keyMap KeyMap
}

func New(groups []ItemGroup) Model {
	return Model{
		groups: groups,
		cursor: 0,
	}
}

func (m Model) SetKeyMap(keyMap KeyMap) Model {
	m.keyMap = keyMap
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.CursorUp):
			m.cursor = max(m.cursor-1, 0)

		case key.Matches(msg, m.keyMap.CursorDown):
			m.cursor = min(m.cursor+1, m.Length()-1)
		}
	}

	return m, nil
}

func (m Model) View() string {
	var (
		idx    = 0
		groups = []string{}
	)

	for _, g := range m.groups {
		items := []string{
			groupNameStyle.Render(g.Name),
		}

		for _, i := range g.Items {
			var style lipgloss.Style
			if m.cursor == idx {
				style = hoverItemStyle
			} else {
				style = lipgloss.NewStyle()
			}

			items = append(items,
				i.Style(style).
					Width(m.width).
					Render(i.Name()),
			)

			idx += 1
		}

		groups = append(groups, lipgloss.JoinVertical(lipgloss.Top, items...))
	}

	return lipgloss.JoinVertical(lipgloss.Top, groups...)
}

func (m Model) Length() int {
	var accu int
	for _, g := range m.groups {
		accu += len(g.Items)
	}

	return accu
}
