package router

import "net/http"

type Multiplexer struct {
	impl *http.ServeMux
}

func New() *Multiplexer {
	mux := http.NewServeMux()
	router := Multiplexer{
		impl: mux,
	}
	return &router
}

func (r *Multiplexer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.impl.ServeHTTP(w, req)
}

func (r *Multiplexer) GET(path string, handler RequestHandler) {
	r.impl.Handle("GET "+path, handler)
}

func (r *Multiplexer) POST(path string, handler RequestHandler) {
	r.impl.Handle("POST "+path, handler)
}

func (r *Multiplexer) DELETE(path string, handler RequestHandler) {
	r.impl.Handle("DELETE "+path, handler)
}

func (r *Multiplexer) PUT(path string, handler RequestHandler) {
	r.impl.Handle("PUT "+path, handler)
}
