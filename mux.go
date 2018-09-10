package mux

import (
	"net/http"
)

type Mux struct {
	routes routes
}

func New() *Mux {
	return &Mux{
		routes: routes{},
	}
}

func (mux *Mux) Handler(req *http.Request) http.Handler {
	h := mux.routes.Match(req.URL.Path)
	if h == nil {
		h = http.NotFoundHandler()
	}
	return h
}

func (mux *Mux) Handle(pattern string, h http.Handler) {
	if h == nil {
		panic("mux: nil handler")
	}

	if err := mux.routes.AddRoute(pattern, h); err != nil {
		panic("mux: duplicate path for `" + pattern + "`")
	}
}

func (mux *Mux) HandleFunc(pattern string, hf http.HandlerFunc) {
	var h http.Handler
	if hf != nil {
		h = http.HandlerFunc(hf)
	}

	mux.Handle(pattern, h)
}

func (mux *Mux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	mux.Handler(req).ServeHTTP(rw, req)
}
