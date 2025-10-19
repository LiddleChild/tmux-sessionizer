package superlist

import (
	"iter"
	"slices"
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/internal/fuzzyfinder"
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

	m.Focus(FocusedModeFilter)
	m.SetItems(groups)

	return m
}

func (m Model) Length() int {
	var length int
	for _, g := range m.GetGroupIter() {
		length += len(g.Items)
	}
	return length
}

func (m Model) SetKeyMap(keyMap KeyMap) Model {
	m.keyMap = keyMap
	return m
}

func (m Model) Focused() bool {
	return m.input.Focused()
}

func (m *Model) Focus(mode FocusedMode) tea.Cmd {
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

func (m Model) GetCursor() int {
	return m.cursor
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

func (m Model) GetGroupIter() iter.Seq2[int, ItemGroup] {
	return func(yield func(int, ItemGroup) bool) {
		for i, group := range m.filteredGroups {
			if len(group.Items) == 0 {
				continue
			}

			if !yield(i, group) {
				return
			}
		}
	}
}

func (m Model) GetItemIter() iter.Seq2[int, Item] {
	return func(yield func(int, Item) bool) {
		var idx int
		for _, group := range m.GetGroupIter() {
			for _, item := range group.Items {
				if filteredItem, ok := item.(*filteredItem); ok {
					item = filteredItem.item
				}

				if !yield(idx, item) {
					return
				}

				idx += 1
			}
		}
	}
}

func (m *Model) SetItems(items []ItemGroup) {
	m.groups = items
	m.filterItems(m.filter.Value())
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
	previewInfo := m.previewList()
	m.ScrollUp(max(0, previewInfo.TopBound-previewInfo.CursorOffset))
	m.ScrollDown(max(0, previewInfo.CursorOffset-previewInfo.BottomBound))
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
		m.listHeight = msg.Height - lipgloss.Height(m.renderFilter())

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
			m.Focus(FocusedModeFilter)

			item := m.GetSelectedItem().(InputItem)

			item.SetValue(m.input.Value())
			return m, SubmitCmd(item.Value(), m.input.Value())

		case m.Focused() && key.Matches(msg, m.keyMap.Cancel):
			m.Focus(FocusedModeFilter)

		case !m.Focused() && key.Matches(msg, m.keyMap.FocusItem):
			return m, m.Focus(FocusedModeItem)
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
	lines := m.renderList()

	previewInfo := m.previewList()
	lines = lines[previewInfo.TopBound:min(len(lines), previewInfo.BottomBound+1)]

	style := lipgloss.NewStyle().
		Height(m.height).
		MaxHeight(m.height).
		Width(m.width).
		MaxWidth(m.width)

	return style.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			m.renderFilter(),
			lipgloss.JoinVertical(lipgloss.Top, lines...),
		),
	)
}
