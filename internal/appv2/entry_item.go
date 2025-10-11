package appv2

import "github.com/charmbracelet/lipgloss"

type entryItem string

func (i entryItem) Name() string {
	return string(i)
}

func (i entryItem) Style(style lipgloss.Style) lipgloss.Style {
	return style
}
