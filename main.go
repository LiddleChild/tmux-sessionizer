package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	AppName = "tmux-sessionpane"
	Version = "v0.2.1"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	m, err := NewModel()
	if err != nil {
		return fmt.Errorf("error initializing app: %w", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	m.program = p

	exitModel, err := p.Run()
	if err != nil {
		return err
	}

	if exitModel.(model).err != nil {
		return err
	}

	return nil
}
