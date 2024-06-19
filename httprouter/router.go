package httprouter

import "net/http"

type Router struct {
	mux *http.ServeMux
}

func New() *Router {
	serverMux := http.NewServeMux()
	router := Router{
		mux: serverMux,
	}
	return &router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) GET(path string, handler RequestHandler) {
	r.mux.Handle("GET "+path, handler)
}

func (r *Router) POST(path string, handler RequestHandler) {
	r.mux.Handle("POST "+path, handler)
}

func (r *Router) DELETE(path string, handler RequestHandler) {
	r.mux.Handle("DELETE "+path, handler)
}

func (r *Router) PUT(path string, handler RequestHandler) {
	r.mux.Handle("PUT "+path, handler)
}
