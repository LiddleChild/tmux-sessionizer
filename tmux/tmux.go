package tmux

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"slices"
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

func AttachSessionCommand(name string) *exec.Cmd {
	if !InTmux() {
		return exec.Command("tmux", "attach-session", "-t", name)
	} else {
		return exec.Command("tmux", "switch-client", "-t", name)
	}
}

// new session

// delete session

// attach session

// rename session
