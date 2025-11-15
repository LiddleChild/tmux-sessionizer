package utils

import (
	"slices"
	"strings"
)

func JoinNonEmpty(s ...string) string {
	s = slices.DeleteFunc(s, func(e string) bool {
		return e == ""
	})
	return strings.Join(s, " ")
}
