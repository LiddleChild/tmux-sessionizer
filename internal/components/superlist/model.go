package superlist

import (
	"slices"
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/internal/fuzzyfinder"
	"github.com/LiddleChild/tmux-sessionpane/internal/utils"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FocusedMode string

const (
	FocusedModeItem   FocusedMode = "Item"
	FocusedModeFilter FocusedMode = "Filter"
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
	groups         []ItemGroup
	filteredGroups []ItemGroup
	cursor         int

	width      int
	height     int
	listHeight int

	yOffset int

	keyMap KeyMap

	input  textinput.Model
	filter textinput.Model
}

func New(groups []ItemGroup) Model {
	input := textinput.New()
	input.Prompt = ""
	input.TextStyle = hoveredItemStyle
	input.PromptStyle = hoveredItemStyle
	input.Cursor.Style = hoveredItemStyle
	input.Cursor.TextStyle = hoveredItemStyle

	filter := textinput.New()
	filter.Prompt = "> "
	filter.PromptStyle = filterPromptStyle
	filter.Cursor.SetMode(cursor.CursorStatic)

	m := Model{
		groups:         []ItemGroup{},
		filteredGroups: []ItemGroup{},
		cursor:         0,
		width:          0,
		height:         0,
		listHeight:     0,
		yOffset:        0,
		keyMap:         KeyMap{},
		input:          input,
		filter:         filter,
	}

	m.SetFocusedMode(FocusedModeFilter)
	m.SetItems(groups)

	return m
}

func (m Model) Length() int {
	var accu int
	for _, g := range m.filteredGroups {
		accu += len(g.Items)
	}

	return accu
}

func (m Model) GetSelectedItem() Item {
	var (
		groups   = m.filteredGroups
		groupIdx = 0
		idx      = m.cursor
	)

	for idx >= len(groups[groupIdx].Items) {
		idx -= len(groups[groupIdx].Items)
		groupIdx += 1
	}

	switch item := groups[groupIdx].Items[idx].(type) {
	case *filteredItem:
		return item.item

	default:
		return item
	}
}

func (m Model) SetKeyMap(keyMap KeyMap) Model {
	m.keyMap = keyMap
	return m
}

func (m Model) Focused() bool {
	return m.input.Focused()
}

func (m *Model) SetFocusedMode(mode FocusedMode) tea.Cmd {
	m.filter.Blur()
	m.input.Blur()

	switch mode {
	case FocusedModeFilter:
		return m.filter.Focus()

	case FocusedModeItem:
		item := m.GetSelectedItem()
		if item, ok := item.(InputItem); ok {
			m.input.SetValue(item.Value())
			m.input.CursorEnd()
			return m.input.Focus()
		}
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

func (m *Model) SetItems(items []ItemGroup) tea.Cmd {
	m.groups = items
	return FilterCmd(m.filter.Value())
}

func (m Model) GetItems() []ItemGroup {
	return m.filteredGroups
}

func (m Model) renderItem(item Item, style lipgloss.Style) string {
	switch item := item.(type) {
	case *filteredItem:
		var (
			builder   strings.Builder
			lastIndex int
		)

		for _, match := range item.matches {
			var (
				start = match.X
				end   = match.Y + 1
			)

			builder.WriteString(style.Render(item.Label()[lastIndex:start]))
			builder.WriteString(style.
				Foreground(lipgloss.Color("3")).
				Render(item.Label()[start:end]))
			lastIndex = end
		}

		builder.WriteString(style.Render(item.Label()[lastIndex:]))

		return builder.String()

	default:
		return item.Label()
	}
}

func (m *Model) filterItems(filter string) {
	if len(filter) == 0 {
		m.filteredGroups = m.groups
		return
	}

	filteredGroups := make([]ItemGroup, 0, len(m.groups))

	for _, group := range m.groups {
		filteredItems := make([]filteredItem, 0, len(group.Items))
		for _, item := range group.Items {
			score, matches := fuzzyfinder.Match(item.Label(), filter)

			if score > 0 {
				filteredItems = append(filteredItems, filteredItem{
					item:    item,
					matches: matches,
					score:   score,
				})
			}
		}

		slices.SortFunc(filteredItems, func(a, b filteredItem) int {
			if a.score == b.score {
				return len(a.item.Label()) - len(b.item.Label())
			} else {
				return b.score - a.score
			}
		})

		items := make([]Item, 0, len(filteredItems))
		for _, fi := range filteredItems {
			items = append(items, &fi)
		}

		group.Items = items
		filteredGroups = append(filteredGroups, group)
	}

	m.filteredGroups = filteredGroups
	m.SetCursor(0)
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

	for _, g := range m.filteredGroups {
		if len(g.Items) == 0 {
			continue
		}

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
	info.BottomBound = utils.Clamp(m.yOffset+m.listHeight-1, 0, currentRenderLen)

	return info
}

func (m Model) render() []string {
	var (
		idx   = 0
		lines = []string{}
	)

	for _, g := range m.filteredGroups {
		if len(g.Items) == 0 {
			continue
		}

		lines = append(lines, groupNameStyle.Render(g.Name))

		for _, i := range g.Items {
			var style lipgloss.Style
			if m.cursor == idx {
				style = hoveredItemStyle
			} else {
				style = lipgloss.NewStyle()
			}

			var itemName string
			if m.Focused() && m.cursor == idx {
				itemName = m.input.View()
			} else {
				itemName = m.renderItem(i, i.Style(style))
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
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case FilterMsg:
		m.filterItems(msg.Value)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.listHeight = msg.Height - 1

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
			m.SetFocusedMode(FocusedModeFilter)

			item := m.GetSelectedItem().(InputItem)

			item.SetValue(m.input.Value())
			return m, SubmitCmd(item.Value(), m.input.Value())

		case m.Focused() && key.Matches(msg, m.keyMap.Cancel):
			m.SetFocusedMode(FocusedModeFilter)

		case !m.Focused() && key.Matches(msg, m.keyMap.FocusItem):
			return m, m.SetFocusedMode(FocusedModeItem)
		}
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	lastFilterValue := m.filter.Value()
	m.filter, cmd = m.filter.Update(msg)
	cmds = append(cmds, cmd)

	if lastFilterValue != m.filter.Value() {
		cmds = append(cmds, FilterCmd(m.filter.Value()))
		lastFilterValue = m.filter.Value()
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	lines := m.render()

	previewInfo := m.preview()
	lines = lines[previewInfo.TopBound:min(len(lines), previewInfo.BottomBound+1)]

	listContent := lipgloss.JoinVertical(lipgloss.Top, lines...)
	listStyle := lipgloss.NewStyle()

	style := lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height).
		Width(m.width).
		MaxWidth(m.width)

	return style.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			m.filter.View(),
			listStyle.Render(listContent),
		),
	)
}
