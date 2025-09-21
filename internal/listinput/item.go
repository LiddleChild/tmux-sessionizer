package listinput

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item interface {
	list.Item

	Label() string
	Value() string
	SetValue(string) tea.Cmd
	Style(lipgloss.Style) lipgloss.Style
}

type item struct {
	Item

	input textinput.Model
}
