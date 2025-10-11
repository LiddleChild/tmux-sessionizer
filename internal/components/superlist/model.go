package superlist

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/log"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	groups []ItemGroup
	cursor int

	width  int
	height int

	keyMap KeyMap

	input textinput.Model
}

func New(items []ItemGroup) Model {
	input := textinput.New()
	input.Prompt = ""
	input.TextStyle = hoveredItemStyle
	input.PromptStyle = hoveredItemStyle
	input.Cursor.Style = hoveredItemStyle
	input.Cursor.TextStyle = hoveredItemStyle

	return Model{
		groups: items,
		cursor: 0,
		width:  0,
		height: 0,
		keyMap: KeyMap{},
		input:  input,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd tea.Cmd
	)

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

		case m.Focused() && key.Matches(msg, m.keyMap.Submit):
			m.input.Blur()

			item := m.GetSelectedItem().(InputItem)

			item.SetValue(m.input.Value())
			return m, SubmitCmd(item.Value(), m.input.Value())

		case m.Focused() && key.Matches(msg, m.keyMap.Cancel):
			m.input.Blur()
		}
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
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
				style = hoveredItemStyle
			} else {
				style = lipgloss.NewStyle()
			}

			itemName := i.Name()
			if m.Focused() && m.cursor == idx {
				itemName = m.input.View()
			}

			items = append(items,
				i.Style(style).
					Width(m.width).
					Render(itemName),
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

func (m Model) GetSelectedItem() Item {
	var (
		groupIdx = 0
		idx      = m.cursor
	)

	for idx >= len(m.groups[groupIdx].Items) {
		idx -= len(m.groups[groupIdx].Items)
		groupIdx += 1
	}

	log.Dump(log.LogLevelDebug, groupIdx)
	log.Dump(log.LogLevelDebug, idx)

	return m.groups[groupIdx].Items[idx]
}

func (m Model) Focused() bool {
	return m.input.Focused()
}

func (m Model) SetKeyMap(keyMap KeyMap) Model {
	m.keyMap = keyMap
	return m
}

func (m *Model) Focus() tea.Cmd {
	item := m.GetSelectedItem()

	if item, ok := item.(InputItem); ok {
		m.input.SetValue(item.Value())
		m.input.CursorEnd()
		return m.input.Focus()
	}

	return nil
}

func (m *Model) SetCursor(cursor int) {
	m.cursor = cursor
}

func (m *Model) SetItems(items []ItemGroup) {
	m.groups = items
}

func (m Model) GetItems() []ItemGroup {
	return m.groups
}
