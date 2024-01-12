package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"code.vin047.com/focus-browser-server/internal/kagi"
	"code.vin047.com/focus-browser-server/internal/search"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var searchClient search.Client

func main() {
	logger := log.New(os.Stdout, "[INFO]: ", log.LstdFlags)
	loggerErr := log.New(os.Stderr, "[ERROR]: ", log.LstdFlags|log.Lshortfile)

	listenPort, exists := os.LookupEnv("LISTEN_PORT")
	if !exists {
		loggerErr.Panicln("LISTEN_PORT environment variable not set")
	}
	kagiApiKey, exists := os.LookupEnv("KAGI_API_KEY")
	if !exists {
		loggerErr.Panicln("KAGI_API_KEY is not set")
	}
	kagiClient := kagi.NewClient(kagiApiKey)
	searchClient = kagiClient

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/ping", func(r chi.Router) {
			r.Get("/", ping)
		})
		r.Route("/v0", func(r chi.Router) {
			r.Post("/search", searchFunc)
		})
	})

	logger.Println("Listening on port", listenPort)
	http.ListenAndServe(":"+listenPort, r)
}

func ping(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "pong")
}

func searchFunc(w http.ResponseWriter, r *http.Request) {
	request := &SearchRequest{}
	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, ErrRenderInvalidRequest(err))
		return
	}

	result, err := searchClient.Search(request.Query)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, RenderSearch(result))
}
