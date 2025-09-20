package listinput

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	list.Model

	delegate itemDelegate
}

func New(items []Item, width int, height int) Model {
	listItems := make([]list.Item, len(items))
	for i, itm := range items {
		input := textinput.New()
		input.TextStyle = lipgloss.NewStyle().Underline(true)

		listItems[i] = item{
			Item:  itm,
			input: input,
		}
	}

	delegate := itemDelegate{}

	return Model{
		Model:    list.New(listItems, delegate, width, height),
		delegate: delegate,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.IsFocused() {
		return m, m.delegate.Update(msg, &m.Model)
	} else {
		var cmd tea.Cmd
		m.Model, cmd = m.Model.Update(msg)
		return m, cmd
	}
}

func (m Model) View() string {
	return m.Model.View()
}

func (m Model) IsFocused() bool {
	isFocused := false
	for _, i := range m.Items() {
		item := i.(item)
		isFocused = isFocused || item.input.Focused()
	}

	return isFocused
}

func (m Model) FocusSelectedItem() tea.Cmd {
	item := m.SelectedItem().(item)

	item.input.SetValue(item.Value())
	item.input.CursorEnd()

	var cmds []tea.Cmd
	cmds = append(cmds, item.input.Focus())
	cmds = append(cmds, m.SetItem(m.Index(), item))

	return tea.Batch(cmds...)
}
