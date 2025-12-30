// Package tmux
package tmux

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/LiddleChild/tmux-sessionizer/internal/log"
)

type session struct {
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
	IsAttached int    `json:"is_attached"`
}

var (
	NoServerRunningErr = errors.New("no server running")
)

func InTmux() bool {
	return os.Getenv("TMUX") != ""
}

func StartServer() error {
	cmd := exec.Command("tmux", "start-server")

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return parseError(err, stderr.String())
	}

	return nil
}

func ListSessions() ([]Session, error) {
	cmd := exec.Command(
		"tmux",
		"list-sessions",
		"-F",
		`{ "name": "#{session_name}", "created_at": #{session_created}, "is_attached": #{session_attached} }`,
	)

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	bs, err := cmd.Output()
	if err != nil {
		return nil, parseError(err, stderr.String())
	}

	decoder := json.NewDecoder(bytes.NewReader(bs))

	sessions := make([]Session, 0)
	for decoder.More() {
		var session session
		if err := decoder.Decode(&session); err != nil {
			return nil, err
		}

		sessions = append(sessions, Session{
			Name:       session.Name,
			CreatedAt:  time.Unix(session.CreatedAt, 0),
			IsAttached: session.IsAttached != 0,
		})
	}

	slices.SortFunc(sessions, func(a, b Session) int {
		return a.CreatedAt.Compare(b.CreatedAt)
	})

	return sessions, nil
}

func AttachSessionCommand(name string) *exec.Cmd {
	if !InTmux() {
		return exec.Command("tmux", "attach-session", "-t", name)
	} else {
		return exec.Command("tmux", "switch-client", "-t", name)
	}
}

func RenameSession(name, newName string) error {
	cmd := exec.Command("tmux", "rename-session", "-t", name, newName)

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return parseError(err, stderr.String())
	}

	return nil
}

func DeleteSession(name string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", name)

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return parseError(err, stderr.String())
	}

	return nil
}

func HasSession(name string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", name)

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		log.Error().Msg(strings.TrimSpace(parseError(err, stderr.String()).Error()))
		return false
	}

	return cmd.ProcessState.ExitCode() == 0
}

func NewDetachedSession(name, workDir string) error {
	cmd := exec.Command("tmux", "new-session", "-d", "-s", name, "-c", workDir)

	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return parseError(err, stderr.String())
	}

	return nil
}

func parseError(err error, stderr string) error {
	if strings.Contains(stderr, "no server running") {
		return fmt.Errorf("%w: %w", err, NoServerRunningErr)
	}

	return fmt.Errorf("%w: %s", err, stderr)
}
