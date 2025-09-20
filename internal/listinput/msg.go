package listinput

import tea "github.com/charmbracelet/bubbletea"

type InputSubmitedMsg struct {
	Index int
	Value string
}

func InputSubmitedCmd(index int, value string) tea.Cmd {
	return func() tea.Msg {
		return InputSubmitedMsg{
			Index: index,
			Value: value,
		}
	}
}
