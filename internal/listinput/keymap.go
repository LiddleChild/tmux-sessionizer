package listinput

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	CursorDown key.Binding
	CursorUp   key.Binding
	Submit     key.Binding
	Cancel     key.Binding
}
