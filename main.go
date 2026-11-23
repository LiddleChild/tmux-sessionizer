package main

import (
	"fmt"

	"github.com/LiddleChild/tmux-sessionizer/internal/app"
	"github.com/LiddleChild/tmux-sessionizer/internal/config"
	"github.com/LiddleChild/tmux-sessionizer/internal/flag"
	"github.com/LiddleChild/tmux-sessionizer/internal/log"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if flag.VersionFlag {
		fmt.Printf("%s %s\n", config.AppName, config.AppVersion)
		return nil
	}

	if err := config.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize: %w", err)
	}

	if err := log.Initialize(flag.DebugFlag); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	m, err := app.New()
	if err != nil {
		return fmt.Errorf("error initializing app: %w", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
