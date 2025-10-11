package appv2

import tea "github.com/charmbracelet/bubbletea"

type ErrMsg error

func ErrCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrMsg(err)
	}
}

type ListTmuxSessionMsg struct{}

func ListTmuxSessionCmd() tea.Msg {
	return ListTmuxSessionMsg{}
}

type SelectAttachedSessionMsg struct{}

func SelectAttachedSessionCmd() tea.Msg {
	return SelectAttachedSessionMsg{}
}
