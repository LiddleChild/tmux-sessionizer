package superlist

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/utils"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type previewInfo struct {
	// first line in bound
	TopBound int

	// last line in bound
	BottomBound int

	// global cursor offset
	CursorOffset int

	// rendering height
	Height int
}

type Model struct {
	groups []ItemGroup
	cursor int

	width  int
	height int

	yOffset int

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
		groups:  items,
		cursor:  0,
		width:   0,
		height:  0,
		yOffset: 0,
		keyMap:  KeyMap{},
		input:   input,
	}
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

func (m *Model) CursorUp() {
	m.cursor = max(m.cursor-1, 0)
}

func (m *Model) CursorDown() {
	m.cursor = min(m.cursor+1, m.Length()-1)
}

func (m *Model) ScrollUp(amount int) {
	m.yOffset -= amount
}

func (m *Model) ScrollDown(amount int) {
	m.yOffset += amount
}

func (m Model) GetCursor() int {
	return m.cursor
}

func (m *Model) SetItems(items []ItemGroup) {
	m.groups = items
}

func (m Model) GetItems() []ItemGroup {
	return m.groups
}

func (m *Model) updateScroll() {
	previewInfo := m.preview()
	m.ScrollUp(max(0, previewInfo.TopBound-previewInfo.CursorOffset))
	m.ScrollDown(max(0, previewInfo.CursorOffset-previewInfo.BottomBound))
}

// TODO: find a way to refactor preview() and render()
func (m Model) preview() previewInfo {
	var (
		info previewInfo

		itemIdx          = 0
		currentRenderLen = 0
	)

	for _, g := range m.groups {
		// group name
		currentRenderLen += 1

		for range g.Items {
			if m.cursor == itemIdx {
				info.CursorOffset = currentRenderLen
			}

			// item
			currentRenderLen += 1
			itemIdx += 1
		}
	}

	info.Height = currentRenderLen
	info.TopBound = utils.Clamp(m.yOffset, 0, currentRenderLen)
	info.BottomBound = utils.Clamp(m.yOffset+m.height-1, 0, currentRenderLen)

	return info
}

func (m Model) render() []string {
	var (
		idx   = 0
		lines = []string{}
	)

	for _, g := range m.groups {
		lines = append(lines, groupNameStyle.Render(g.Name))

		for _, i := range g.Items {
			var style lipgloss.Style
			if m.cursor == idx {
				style = hoveredItemStyle
			} else {
				style = lipgloss.NewStyle()
			}

			itemName := i.Label()
			if m.Focused() && m.cursor == idx {
				itemName = m.input.View()
			}

			lines = append(lines,
				i.Style(style).
					Width(m.width).
					Render(itemName),
			)

			idx += 1
		}
	}

	return lines
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

		m.updateScroll()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.CursorUp):
			m.CursorUp()
			m.updateScroll()

		case key.Matches(msg, m.keyMap.CursorDown):
			m.CursorDown()
			m.updateScroll()

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
	lines := m.render()

	previewInfo := m.preview()
	lines = lines[previewInfo.TopBound:min(len(lines), previewInfo.BottomBound+1)]

	return lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height).
		Width(m.width).
		MaxWidth(m.width).
		Render(
			lipgloss.JoinVertical(lipgloss.Top, lines...),
		)
}
