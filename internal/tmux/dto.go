package tmux

import "time"

type Session struct {
	Name       string
	CreatedAt  time.Time
	IsAttached bool
}
