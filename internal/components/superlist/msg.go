package superlist

import tea "github.com/charmbracelet/bubbletea"

type SubmitMsg struct {
	OldValue string
	NewValue string
}

func SubmitCmd(oldValue, newValue string) tea.Cmd {
	return func() tea.Msg {
		return SubmitMsg{
			OldValue: oldValue,
			NewValue: newValue,
		}
	}
}
