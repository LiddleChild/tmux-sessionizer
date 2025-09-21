// Package log
package log

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/LiddleChild/tmux-sessionpane/internal/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

const (
	DebugEntry = "debug.log"
	ErrorEntry = "error.log"
)

var (
	DebugEntryPath = path.Join(config.BaseConfigPath, "debug.log")
	ErrorEntryPath = path.Join(config.BaseConfigPath, "error.log")
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelError
)

var (
	DebugFlag = flag.Bool("debug", false, "debug")

	logLevel LogLevel
	entry    *os.File
)

func Init() error {
	flag.Parse()

	var err error
	if *DebugFlag {
		entry, err = os.OpenFile(DebugEntryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		logLevel = LogLevelDebug
	} else {
		entry, err = os.OpenFile(ErrorEntryPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		logLevel = LogLevelInfo
	}

	return err
}

func printTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func printLogLevel(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("13")).
			Render("DEBUG")

	case LogLevelInfo:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Render("INFO")

	case LogLevelError:
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("1")).Bold(true).
			Render("ERROR")

	default:
		return ""
	}
}

func print(s ...string) {
	fmt.Fprint(entry, strings.Join(s, " "))
}

func Printlnf(level LogLevel, format string, a ...any) {
	if entry == nil {
		return
	}

	if level >= logLevel {
		print(
			printTimestamp(),
			printLogLevel(level),
			fmt.Sprintf(format, a...),
			"\n",
		)
	}
}

func Dump(level LogLevel, v any) {
	if entry == nil {
		return
	}

	if level >= logLevel {
		print(
			printTimestamp(),
			printLogLevel(level),
			spew.Sdump(v),
		)
	}
}
