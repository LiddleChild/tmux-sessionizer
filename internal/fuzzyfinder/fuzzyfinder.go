package fuzzyfinder

type Source interface {
	Get(int) string
	Len() int
}

type Match struct {
	Index          int
	MatchedIndices []int
	Score          int
}

type Algorithm interface {
	Find(Source, string) []Match
}
