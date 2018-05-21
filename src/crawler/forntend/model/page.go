package model

type SearchResult struct {
	Hint int64
	Start int
	Query string
	PreFrom int
	NextFrom int
	Items []interface{}
}
