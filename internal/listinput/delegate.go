package listinput

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ list.ItemDelegate = (*itemDelegate)(nil)

type itemDelegate struct {
	keyMap KeyMap
}

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	for i, listItem := range m.Items() {
		listInputItem := listItem.(item)

		if listInputItem.input.Focused() {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch {
				case key.Matches(msg, d.keyMap.Submit):
					cmds = append(cmds, listInputItem.SetValue(listInputItem.input.Value()))
					fallthrough

				case key.Matches(msg, d.keyMap.Cancel):
					listInputItem.input.Blur()
				}
			}
		}

		var cmd tea.Cmd
		listInputItem.input, cmd = listInputItem.input.Update(msg)
		cmds = append(cmds, cmd)

		cmd = m.SetItem(i, list.Item(listInputItem))
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(item)
	if !ok {
		return
	}

	str := item.Label()
	str += strings.Repeat(" ", max(m.Width()-len(str), 0))

	if index == m.Index() {
		style := item.Style(selectedItemStyle)

		if item.input.Focused() {
			item.input.Width = m.Width()
			str = style.Render(item.input.View())
		} else {
			str = style.Render(str)
		}
	} else {
		str = item.Style(lipgloss.NewStyle()).Render(str)
	}

	fmt.Fprint(w, str)
}
