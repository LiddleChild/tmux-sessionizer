package superlist

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	Submit     key.Binding
	Cancel     key.Binding
}
