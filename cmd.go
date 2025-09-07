package main

import tea "github.com/charmbracelet/bubbletea"

type AttachSessionMsg struct {
	Name string
}

func AttachSession(name string) tea.Cmd {
	return func() tea.Msg {
		return AttachSessionMsg{Name: name}
	}
}
