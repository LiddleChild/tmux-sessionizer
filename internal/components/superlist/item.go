package superlist

import "github.com/charmbracelet/lipgloss"

type Item interface {
	Name() string
	Style(lipgloss.Style) lipgloss.Style
}

type ItemGroup struct {
	Name  string
	Items []Item
}
