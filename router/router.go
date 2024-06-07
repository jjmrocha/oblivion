package router

import "net/http"

type Router struct {
	multiplexer *http.ServeMux
}

func New() *Router {
	mux := http.NewServeMux()
	router := Router{
		multiplexer: mux,
	}
	return &router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.multiplexer.ServeHTTP(w, req)
}

func (r *Router) GET(path string, handler RequestHandler) {
	r.multiplexer.Handle("GET "+path, &handlerWrapper{handler})
}

func (r *Router) POST(path string, handler RequestHandler) {
	r.multiplexer.Handle("POST "+path, &handlerWrapper{handler})
}

func (r *Router) DELETE(path string, handler RequestHandler) {
	r.multiplexer.Handle("DELETE "+path, &handlerWrapper{handler})
}

func (r *Router) PUT(path string, handler RequestHandler) {
	r.multiplexer.Handle("PUT "+path, &handlerWrapper{handler})
}
