package main

import (
	"fmt"

	"github.com/LiddleChild/tmux-sessionpane/internal/config"
	"github.com/LiddleChild/tmux-sessionpane/internal/log"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := config.Init(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	if err := log.Init(); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	m, err := NewModel()
	if err != nil {
		return fmt.Errorf("error initializing app: %w", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
