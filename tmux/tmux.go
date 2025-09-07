package tmux

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"slices"
	"syscall"
	"time"
)

type session struct {
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
	IsAttached int    `json:"is_attached"`
}

func InTmux() bool {
	return os.Getenv("TMUX") != ""
}

func ListSession() ([]Session, error) {
	bs, err := exec.Command(
		"tmux",
		"list-sessions",
		"-F",
		`{ "name": "#{session_name}", "created_at": #{session_created}, "is_attached": #{session_attached} }`,
	).Output()
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(bs))

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

func AttachSession(name string) error {
	tmux, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}

	var args []string
	if !InTmux() {
		args = []string{tmux, "attach-session", "-t", name}
	} else {
		args = []string{tmux, "switch-client", "-t", name}
	}

	return syscall.Exec(tmux, args, os.Environ())
}

// new session

// delete session

// attach session

// rename session
