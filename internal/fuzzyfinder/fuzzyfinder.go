package fuzzyfinder

import (
	"strings"

	"github.com/LiddleChild/tmux-sessionpane/internal/types"
)

func Match(str, pattern string) (int, []types.Pair[int]) {
	return smithWaterman(str, pattern)
}

func smithWaterman(str, pattern string) (int, []types.Pair[int]) {
	const (
		Matched    = 5
		Mismatched = -3
		Gap        = -1
	)

	var (
		strLen     = len(str)
		patternLen = len(pattern)
	)

	str = strings.ToLower(str)
	pattern = strings.ToLower(pattern)

	scores := make([][]int, strLen+1)
	for i := range strLen + 1 {
		scores[i] = make([]int, patternLen+1)
	}

	var (
		score = 0
		mxi   = 0
		mxj   = 0
	)

	for i := 1; i <= strLen; i++ {
		for j := 1; j <= patternLen; j++ {
			var (
				top      = max(0, scores[i-1][j]+Gap)
				left     = max(0, scores[i][j-1]+Gap)
				diagonal = scores[i-1][j-1]
			)

			if str[i-1] == pattern[j-1] {
				diagonal += Matched
			} else {
				diagonal += Mismatched
			}

			scores[i][j] = max(diagonal, max(left, top))

			if scores[i][j] > score {
				score = scores[i][j]
				mxi = i
				mxj = j
			}
		}
	}

	indices := []int{}
	for scores[mxi][mxj] > 0 {
		if str[mxi-1] == pattern[mxj-1] {
			indices = append([]int{mxi - 1}, indices...)
		}

		localMx := max(scores[mxi-1][mxj-1], max(scores[mxi-1][mxj], scores[mxi][mxj-1]))
		if scores[mxi-1][mxj-1] == localMx {
			mxi -= 1
			mxj -= 1
		} else if scores[mxi-1][mxj] == localMx {
			mxi -= 1
		} else if scores[mxi][mxj-1] == localMx {
			mxj -= 1
		}
	}

	if score*2 < patternLen*Matched {
		return -1, []types.Pair[int]{}
	}

	return score, toRange(indices)
}

func toRange(indices []int) []types.Pair[int] {
	var (
		ranges = []types.Pair[int]{}
		i      = 0
		j      = 0
	)

	for i < len(indices) && j < len(indices) {
		for j+1 < len(indices) {
			if indices[j+1]-indices[j] > 1 {
				break
			}

			j += 1
		}

		ranges = append(ranges, types.Pair[int]{
			X: indices[i],
			Y: indices[j],
		})

		i = j + 1
		j += 1
	}

	return ranges
}
