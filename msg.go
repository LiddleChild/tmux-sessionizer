package main

import tea "github.com/charmbracelet/bubbletea"

type QuitWithErrMsg struct {
	err error
}

func QuitWithErr(err error) tea.Cmd {
	return func() tea.Msg {
		return QuitWithErrMsg{err: err}
	}
}

type ListTmuxSessionMsg struct{}

func ListTmuxSessionCmd() tea.Msg {
	return ListTmuxSessionMsg{}
}
