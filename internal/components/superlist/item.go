package superlist

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/types"
	"github.com/charmbracelet/lipgloss"
)

type InputItem interface {
	Value() string
	SetValue(string)

	Item
}

type Item interface {
	Label() string
	Style(lipgloss.Style) lipgloss.Style
}

type ItemGroup struct {
	Name  string
	Items []Item
}

var _ Item = (*filteredItem)(nil)

type filteredItem struct {
	item    Item
	score   int
	matches []types.Pair[int]
}

func (i filteredItem) Label() string {
	return i.item.Label()
}

func (i filteredItem) Style(style lipgloss.Style) lipgloss.Style {
	return i.item.Style(style)
}
