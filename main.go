package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	AppName = "tmux-sessionpane"
	Version = "v0.2.1"
)

var (
	debugFlag = flag.Bool("debug", false, "debug")
)

func main() {
	flag.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var dump *os.File
	if *debugFlag {
		var err error
		dump, err = os.OpenFile("messages.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			os.Exit(1)
		}

		fmt.Fprintln(dump, "DEBUG MODE")
	}

	m, err := NewModel(dump)
	if err != nil {
		return fmt.Errorf("error initializing app: %w", err)
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	exitModel, err := p.Run()
	if err != nil {
		return err
	}

	if exitModel.(model).err != nil {
		return err
	}

	return nil
}
