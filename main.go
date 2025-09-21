package main

import (
	"fmt"
	"os"

	"github.com/LiddleChild/tmux-sessionpane/internal/log"
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
	if *log.DebugFlag {
		var err error
		log.LogFile, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}

		fmt.Fprintln(log.LogFile, "DEBUG MODE")
	}

	m, err := NewModel(log.LogFile)
	if err != nil {
		return fmt.Errorf("error initializing app: %w", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
