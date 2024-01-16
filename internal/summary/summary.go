package summary

type Client interface {
	Summarise(url string) (*Result, error)
}

type Result struct {
	Response string
}
