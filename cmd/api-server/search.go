package main

import (
	"net/http"

	"code.vin047.com/focus-browser-server/internal/search"
	"github.com/go-chi/render"
)

type SearchRequest struct {
	Query string `json:"query"`
}

func (s *SearchRequest) Bind(r *http.Request) error {
	return nil
}

type SearchResponse struct {
	HTTPStatusCode int                       `json:"-"` // http response status code
	Response       string                    `json:"response"`
	References     []SearchResponseReference `json:"references"`
}

type SearchResponseReference struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func (s *SearchResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, s.HTTPStatusCode)
	return nil
}

func RenderSearch(result *search.Result) render.Renderer {
	return &SearchResponse{
		HTTPStatusCode: http.StatusOK,
		Response:       result.Response,
		References: func() []SearchResponseReference {
			var references []SearchResponseReference
			for _, reference := range result.References {
				references = append(references, SearchResponseReference{
					Title: reference.Title,
					Url:   reference.Url,
				})
			}
			return references
		}(),
	}
}
