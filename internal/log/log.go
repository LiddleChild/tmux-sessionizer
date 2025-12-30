// Package log
package log

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/LiddleChild/tmux-sessionizer/internal/colors"
	"github.com/LiddleChild/tmux-sessionizer/internal/config"
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

type LogBuilder struct {
	Level LogLevel
}

var (
	DebugFlag bool

	logLevel LogLevel
	entry    *os.File
)

func Init() error {
	flag.BoolVar(&DebugFlag, "debug", false, "debug")

	flag.Parse()

	var err error
	if DebugFlag {
		entry, err = os.OpenFile(DebugEntryPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		logLevel = LogLevelDebug
	} else {
		entry, err = os.OpenFile(ErrorEntryPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		logLevel = LogLevelInfo
	}

	return err
}

func Info() LogBuilder {
	return LogBuilder{
		Level: LogLevelInfo,
	}
}

func Error() LogBuilder {
	return LogBuilder{
		Level: LogLevelError,
	}
}

func Debug() LogBuilder {
	return LogBuilder{
		Level: LogLevelDebug,
	}
}

func (b LogBuilder) Dump(v any) {
	print(
		b.Level,
		printTimestamp(),
		printLogLevel(b.Level),
		spew.Sdump(v),
	)
}

func (b LogBuilder) Msg(s string) {
	print(
		b.Level,
		printTimestamp(),
		printLogLevel(b.Level),
		s,
		"\n",
	)
}

func printTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func printLogLevel(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return lipgloss.NewStyle().
			Foreground(colors.BrightMagenta).
			Render("DEBUG")

	case LogLevelInfo:
		return lipgloss.NewStyle().
			Foreground(colors.Green).
			Render("INFO")

	case LogLevelError:
		return lipgloss.NewStyle().
			Foreground(colors.Red).Bold(true).
			Render("ERROR")

	default:
		return ""
	}
}

func print(level LogLevel, s ...string) {
	if level >= logLevel {
		fmt.Fprint(entry, strings.Join(s, " "))
	}
}
