// Package listinput
package listinput

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list     list.Model
	delegate itemDelegate
}

func New(items []Item) Model {
	listItems := make([]list.Item, len(items))
	for i, listItem := range items {
		input := textinput.New()
		input.Prompt = ""
		input.TextStyle = selectedItemStyle
		input.PromptStyle = selectedItemStyle
		input.Cursor.Style = selectedItemStyle
		input.Cursor.TextStyle = selectedItemStyle

		listItems[i] = item{
			Item:  listItem,
			input: input,
		}
	}

	delegate := itemDelegate{}

	list := list.New(listItems, delegate, 0, len(listItems))
	list.SetFilteringEnabled(false)
	list.SetShowStatusBar(false)
	list.SetShowTitle(false)
	list.SetShowHelp(false)
	list.SetShowPagination(false)

	return Model{
		list:     list,
		delegate: delegate,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.IsFocused() {
		return m, m.delegate.Update(msg, &m.list)
	} else {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	return m.list.View()
}

func (m Model) IsFocused() bool {
	isFocused := false
	for _, i := range m.list.Items() {
		item := i.(item)
		isFocused = isFocused || item.input.Focused()
	}

	return isFocused
}

func (m Model) FocusSelectedItem() tea.Cmd {
	item := m.list.SelectedItem().(item)

	item.input.SetValue(item.Value())
	item.input.CursorEnd()

	var cmds []tea.Cmd
	cmds = append(cmds, item.input.Focus())
	cmds = append(cmds, m.list.SetItem(m.list.Index(), item))

	return tea.Batch(cmds...)
}

func (m Model) Index() int {
	return m.list.Index()
}

func (m Model) SelectedItem() Item {
	return m.list.SelectedItem().(item).Item
}

func (m Model) Items() []Item {
	items := make([]Item, len(m.list.Items()))
	for i, item := range m.list.Items() {
		items[i] = item.(Item)
	}

	return items
}

func (m *Model) CursorUp() {
	m.list.CursorUp()
}

func (m *Model) SetWidth(width int) {
	m.list.SetWidth(width)
}

func (m *Model) SetKeyMap(keyMap KeyMap) {
	m.delegate.keyMap = keyMap

	m.list.SetDelegate(m.delegate)

	m.list.KeyMap = list.KeyMap{
		CursorUp:   keyMap.CursorUp,
		CursorDown: keyMap.CursorDown,
	}
}

func (m *Model) SetItems(items []Item) tea.Cmd {
	listItems := make([]list.Item, len(items))
	for i, listItem := range items {
		input := textinput.New()
		input.Prompt = ""
		input.TextStyle = selectedItemStyle
		input.PromptStyle = selectedItemStyle
		input.Cursor.Style = selectedItemStyle
		input.Cursor.TextStyle = selectedItemStyle

		listItems[i] = item{
			Item:  listItem,
			input: input,
		}
	}

	m.list.SetHeight(len(listItems) + 1)

	return m.list.SetItems(listItems)
}
