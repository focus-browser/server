package main

import (
	"net/http"

	"code.vin047.com/focus-browser-server/internal/summary"
	"github.com/go-chi/render"
)

type SummaryRequest struct {
	Url string `json:"url"`
}

func (s *SummaryRequest) Bind(r *http.Request) error {
	return nil
}

type SummaryResponse struct {
	HTTPStatusCode int    `json:"-"` // http response status code
	Response       string `json:"response"`
}

func (s *SummaryResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, s.HTTPStatusCode)
	return nil
}

func RenderSummary(result *summary.Result) render.Renderer {
	return &SummaryResponse{
		HTTPStatusCode: http.StatusOK,
		Response:       result.Response,
	}
}
