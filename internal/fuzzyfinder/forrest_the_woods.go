// a modified version of github.com/sahilm/fuzzy
// a copy-and-paste with slight change
package fuzzyfinder

import (
	"math"
	"slices"
	"unicode"
	"unicode/utf8"
)

type partialMatch struct {
	MatchedIndices []int
	Score          int
}

type ForrestTheWoods struct {
	separators                     []rune
	firstCharMatchBonus            int
	matchFollowingSeparatorBonus   int
	camelCaseMatchBonus            int
	adjacentMatchBonus             int
	unmatchedLeadingCharPenalty    int
	maxUnmatchedLeadingCharPenalty int
}

func NewForrestTheWoods() Algorithm {
	return &ForrestTheWoods{
		separators:                     []rune("/-_ .\\"),
		firstCharMatchBonus:            10,
		matchFollowingSeparatorBonus:   20,
		camelCaseMatchBonus:            20,
		adjacentMatchBonus:             5,
		unmatchedLeadingCharPenalty:    -5,
		maxUnmatchedLeadingCharPenalty: -15,
	}
}

func (a *ForrestTheWoods) Find(src Source, pattern string) []Match {
	if len(pattern) == 0 {
		return []Match{}
	}

	matches := make([]Match, 0)
	for i := 0; i < src.Len(); i++ {
		match := a.bestMatchingSingleSource(src.Get(i), pattern)
		if len(match.MatchedIndices) == 0 {
			continue
		}

		matches = append(matches, Match{
			Index:          i,
			MatchedIndices: match.MatchedIndices,
			Score:          match.Score,
		})
	}

	slices.SortStableFunc(matches, func(a, b Match) int {
		return b.Score - a.Score
	})

	return matches
}

func (a *ForrestTheWoods) bestMatchingSingleSource(src, pattern string) partialMatch {
	var bestMatch partialMatch
	bestMatch.Score = math.MinInt

	for {
		match := a.matchSingleSource(src, pattern)
		if len(match.MatchedIndices) == 0 {
			break
		}

		if match.Score > bestMatch.Score {
			bestMatch = match
		}

		src = a.replaceRuneIndex(src, match.MatchedIndices[0], rune(' '))
	}

	return bestMatch
}

func (a *ForrestTheWoods) matchSingleSource(src, pattern string) partialMatch {
	runes := []rune(pattern)

	var match partialMatch

	var (
		score                  = 0
		bestScore              = -1
		currAdjacentMatchBonus = 0

		patternIndex = 0
		matchedIndex = -1
	)

	var (
		last      rune
		lastIndex int
	)

	var (
		candidate     rune
		candidateSize int
	)

	nextc, nextSize := utf8.DecodeRuneInString(src)
	for j := 0; j < len(src); j += candidateSize {
		candidate, candidateSize = nextc, nextSize
		if a.equalFold(candidate, runes[patternIndex]) {
			score = 0
			if j == 0 {
				score += a.firstCharMatchBonus
			}
			if unicode.IsLower(last) && unicode.IsUpper(candidate) {
				score += a.camelCaseMatchBonus
			}
			if j != 0 && a.isSeparator(last) {
				score += a.matchFollowingSeparatorBonus
			}
			if len(match.MatchedIndices) > 0 {
				lastMatch := match.MatchedIndices[len(match.MatchedIndices)-1]
				bonus := a.adjacentCharBonus(lastIndex, lastMatch, currAdjacentMatchBonus)
				score += bonus
				// adjacent matches are incremental and keep increasing based on previous adjacent matches
				// thus we need to maintain the current match bonus
				currAdjacentMatchBonus += bonus
			}
			if score > bestScore {
				bestScore = score
				matchedIndex = j
			}
		}
		var nextp rune
		if patternIndex < len(runes)-1 {
			nextp = runes[patternIndex+1]
		}
		if j+candidateSize < len(src) {
			if src[j+candidateSize] < utf8.RuneSelf { // Fast path for ASCII
				nextc, nextSize = rune(src[j+candidateSize]), 1
			} else {
				nextc, nextSize = utf8.DecodeRuneInString(src[j+candidateSize:])
			}
		} else {
			nextc, nextSize = 0, 0
		}
		// We apply the best score when we have the next match coming up or when the search string has ended.
		// Tracking when the next match is coming up allows us to exhaustively find the best match and not necessarily
		// the first match.
		// For example given the pattern "tk" and search string "The Black Knight", exhaustively matching allows us
		// to match the second k thus giving this string a higher score.
		if a.equalFold(nextp, nextc) || nextc == 0 {
			if matchedIndex > -1 {
				if len(match.MatchedIndices) == 0 {
					penalty := matchedIndex * a.unmatchedLeadingCharPenalty
					bestScore += max(penalty, a.maxUnmatchedLeadingCharPenalty)
				}
				match.Score += bestScore
				match.MatchedIndices = append(match.MatchedIndices, matchedIndex)
				score = 0
				bestScore = -1
				patternIndex++
			}
		}
		lastIndex = j
		last = candidate
	}
	// apply penalty for each unmatched character
	penalty := len(match.MatchedIndices) - len(src)
	match.Score += penalty
	if len(match.MatchedIndices) != len(runes) {
		return partialMatch{}
	}

	return match
}

// Taken from strings.EqualFold
func (a *ForrestTheWoods) equalFold(tr, sr rune) bool {
	if tr == sr {
		return true
	}
	if tr < sr {
		tr, sr = sr, tr
	}
	// Fast check for ASCII.
	if tr < utf8.RuneSelf {
		// ASCII, and sr is upper case.  tr must be lower case.
		if 'A' <= sr && sr <= 'Z' && tr == sr+'a'-'A' {
			return true
		}
		return false
	}

	// General case. SimpleFold(x) returns the next equivalent rune > x
	// or wraps around to smaller values.
	r := unicode.SimpleFold(sr)
	for r != sr && r < tr {
		r = unicode.SimpleFold(r)
	}
	return r == tr
}

func (a *ForrestTheWoods) adjacentCharBonus(i int, lastMatch int, currentBonus int) int {
	if lastMatch == i {
		return currentBonus*2 + a.adjacentMatchBonus
	}
	return 0
}

func (a *ForrestTheWoods) isSeparator(r rune) bool {
	return slices.Contains(a.separators, r)
}

func (a *ForrestTheWoods) replaceRuneIndex(s string, idx int, replace rune) string {
	rs := []rune(s)
	for i := range rs {
		if i == idx {
			rs[i] = replace
		}
	}
	return string(rs)
}
