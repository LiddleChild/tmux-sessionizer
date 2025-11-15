package superlist

import (
	"github.com/LiddleChild/tmux-sessionizer/internal/fuzzyfinder"
	"github.com/charmbracelet/lipgloss"
)

type InputItem interface {
	Value() string
	SetValue(string)

	Item
}

type Item interface {
	Label() string
	Suffix() string
	Style(lipgloss.Style) lipgloss.Style
}

var _ fuzzyfinder.Source = (*ItemGroup)(nil)

type ItemGroup struct {
	Name  string
	Items []Item
}

func (g ItemGroup) Get(i int) string {
	return g.Items[i].Label()
}

func (g ItemGroup) Len() int {
	return len(g.Items)
}
