package superlist

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/types"
	"github.com/charmbracelet/lipgloss"
)

var _ Item = (*filteredItem)(nil)

type filteredItem struct {
	item    Item
	score   int
	matches []types.Pair[int]
}

func (item filteredItem) Label() string {
	return item.item.Label()
}

func (item filteredItem) Style(style lipgloss.Style) lipgloss.Style {
	return style
}
