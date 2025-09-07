package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/tmux"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ list.Item = (*item)(nil)

type item tmux.Session

func (i item) Title() string {
	if i.IsAttached {
		return i.Name + " (attached)"
	}

	return i.Name
}

func (i item) FilterValue() string {
	return i.Name
}

var _ list.ItemDelegate = (*itemDelegate)(nil)

var (
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Background(lipgloss.Color("8")).Bold(true)
)

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	session, ok := listItem.(item)
	if !ok {
		return
	}

	str := session.Title()
	str += strings.Repeat(" ", max(m.Width()-len(str), 0))

	if index == m.Index() {
		str = selectedItemStyle.Render(str)
	}

	fmt.Fprint(w, str)
}

var _ tea.Model = (*model)(nil)

type model struct {
	program *tea.Program
	err     error

	list list.Model
}

func NewModel() (*model, error) {
	sessions, err := tmux.ListSession()
	if err != nil {
		return nil, err
	}

	items := []list.Item{}
	for _, session := range sessions {
		items = append(items, item(session))
	}

	l := list.New(items, itemDelegate{}, 0, len(sessions)+2)
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)

	return &model{list: l}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			session := m.list.SelectedItem().(item)

			execProcessCmd := tea.ExecProcess(tmux.AttachSessionCommand(session.Name), func(err error) tea.Msg {
				return QuitWithErr(err)
			})

			return m, tea.Sequence(
				execProcessCmd,
				tea.Quit,
			)
		}

	case QuitWithErrMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("%s %s", AppName, Version))
	builder.WriteByte('\n')
	builder.WriteString(m.list.Help.View(m.list))
	builder.WriteByte('\n')
	builder.WriteByte('\n')
	builder.WriteString(m.list.View())

	return builder.String()
}
