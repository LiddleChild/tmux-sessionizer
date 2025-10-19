package fuzzyfinder

import (
	"slices"
	"strings"
)

type scoredItem struct {
	index          int
	score          int
	matchedIndices []int
	len            int
}

var _ Algorithm = (*SmithWaterman)(nil)

type SmithWaterman struct {
	matched    int
	mismatched int
	gap        int
}

func NewSmithWaterman() Algorithm {
	return &SmithWaterman{
		matched:    5,
		mismatched: -3,
		gap:        -1,
	}
}

func (sw *SmithWaterman) Find(src Source, pattern string) []Match {
	scoredItems := make([]scoredItem, 0, src.Len())
	for i := 0; i < src.Len(); i += 1 {
		score, indices := sw.score(src.Get(i), pattern)

		if score*2 >= len(pattern)*sw.matched {
			scoredItems = append(scoredItems, scoredItem{
				index:          i,
				score:          score,
				matchedIndices: indices,
				len:            len(src.Get(i)),
			})
		}
	}

	slices.SortStableFunc(scoredItems, func(a, b scoredItem) int {
		if a.score == b.score {
			return a.len - b.len
		} else {
			return b.score - a.score
		}
	})

	matches := make([]Match, 0, len(scoredItems))
	for _, scoredItem := range scoredItems {
		matches = append(matches, Match{
			Index:          scoredItem.index,
			MatchedIndices: scoredItem.matchedIndices,
			Score:          scoredItem.score,
		})
	}

	return matches
}

func (sw *SmithWaterman) score(str, pattern string) (int, []int) {
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
				top      = max(0, scores[i-1][j]+sw.gap)
				left     = max(0, scores[i][j-1]+sw.gap)
				diagonal = scores[i-1][j-1]
			)

			if str[i-1] == pattern[j-1] {
				diagonal += sw.matched
			} else {
				diagonal += sw.mismatched
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

	return score, indices
}
