package appv2

import (
	"github.com/LiddleChild/tmux-sessionpane/internal/utils"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

var _ help.KeyMap = (*KeyMap)(nil)

func (k KeyMap) ShortHelp() []key.Binding { return []key.Binding{} }
func (k KeyMap) FullHelp() [][]key.Binding {
	return utils.Transpose(
		[][]key.Binding{
			{k.Up, k.Down, k.Quit},
			{k.Select, k.Rename, k.Delete},
		},
	)
}

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Quit   key.Binding
	Select key.Binding
	Rename key.Binding
	Delete key.Binding
}

var keyMap = KeyMap{
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

var _ help.KeyMap = (*FocusedKeyMap)(nil)

func (k FocusedKeyMap) ShortHelp() []key.Binding { return []key.Binding{} }
func (k FocusedKeyMap) FullHelp() [][]key.Binding {
	return utils.Transpose(
		[][]key.Binding{
			{k.Submit, k.Cancel},
		},
	)
}

type FocusedKeyMap struct {
	Submit key.Binding
	Cancel key.Binding
}

var focusedKeyMap = FocusedKeyMap{
	Submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "rename"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}
