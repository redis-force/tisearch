package model

//Field
type Field struct {
	Name  string
	Value string
}

type SearchResult struct {
	Rows []Row
}

type Row struct {
	DocID  int64
	Render map[string][]string
}
