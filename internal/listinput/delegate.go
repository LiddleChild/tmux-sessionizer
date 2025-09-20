package listinput

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ list.ItemDelegate = (*itemDelegate)(nil)

var (
	selectedItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("5")).
		Background(lipgloss.Color("8")).
		Bold(true)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int { return 1 }

func (d itemDelegate) Spacing() int { return 0 }

func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd
	for i, listItem := range m.Items() {
		listInputItem := listItem.(item)

		if listInputItem.input.Focused() {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "enter", "esc":
					listInputItem.input.Blur()
					cmds = append(cmds, listInputItem.OnValueChange(listInputItem.input.Value()))
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

	str := item.Name()
	str += strings.Repeat(" ", max(m.Width()-len(str), 0))

	if index == m.Index() {
		if item.input.Focused() {
			str = item.input.View()
		} else {
			str = selectedItemStyle.Render(str)
		}
	}

	fmt.Fprint(w, str)
}
