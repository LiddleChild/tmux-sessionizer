package fuzzyfinder

import (
	"github.com/sahilm/fuzzy"
)

type sourceAdapter struct {
	Source
}

func (sa sourceAdapter) String(i int) string {
	return sa.Source.Get(i)
}

func (sa sourceAdapter) Len() int {
	return sa.Source.Len()
}

type ForrestTheWoods struct{}

func NewForrestTheWoods() Algorithm {
	return &ForrestTheWoods{}
}

func (s *ForrestTheWoods) Find(src Source, pattern string) []Match {
	var (
		results = fuzzy.FindFrom(pattern, sourceAdapter{src})
		matches = make([]Match, 0, len(results))
	)

	for _, match := range results {
		matches = append(matches, Match{
			Index:          match.Index,
			MatchedIndices: match.MatchedIndexes,
			Score:          match.Score,
		})
	}

	return matches
}
