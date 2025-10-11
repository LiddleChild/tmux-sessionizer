package superlist

import (
	"github.com/charmbracelet/lipgloss"
)

type InputItem interface {
	Value() string
	SetValue(string)

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
