package listinput

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Item interface {
	list.Item

	Label() string
	Value() string
	SetValue(val string) tea.Cmd
}

type item struct {
	Item

	input textinput.Model
}
