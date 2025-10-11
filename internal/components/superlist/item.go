package superlist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputItem interface {
	Value() string
	SetValue(string) tea.Cmd

	Item
}

type Item interface {
	Name() string
	Style(lipgloss.Style) lipgloss.Style
}

type ItemGroup struct {
	Name  string
	Items []Item
}
