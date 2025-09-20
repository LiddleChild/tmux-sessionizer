package listinput

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

type Item interface {
	list.Item

	Label() string
	Value() string
}

type item struct {
	Item

	input textinput.Model
}
