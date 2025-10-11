package appv2

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*keyMap)(nil)

func (k keyMap) ShortHelp() []key.Binding { return []key.Binding{} }
func (k keyMap) FullHelp() [][]key.Binding {
	return utils.Transpose(
		[][]key.Binding{
			{k.Up, k.Down, k.Quit},
			{k.Select, k.Rename, k.Delete},
		},
	)
}

type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Select key.Binding
	Rename key.Binding
	Delete key.Binding
}

var keymap = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑ / k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓ / j", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q / esc / ctrl+c", "quit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "goto"),
	),
	Rename: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rename"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}

var _ help.KeyMap = (*focusedKeyMap)(nil)

func (k focusedKeyMap) ShortHelp() []key.Binding { return []key.Binding{} }
func (k focusedKeyMap) FullHelp() [][]key.Binding {
	return utils.Transpose(
		[][]key.Binding{
			{k.Submit, k.Cancel},
		},
	)
}

type focusedKeyMap struct {
	Submit key.Binding
	Cancel key.Binding
}

var focusedKeymap = focusedKeyMap{
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "rename"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}
