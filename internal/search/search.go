package search

type Client interface {
	Search(query string) (*Result, error)
}

type Result struct {
	Response   string
	References []Reference
}

type Reference struct {
	Title string
	Url   string
}
