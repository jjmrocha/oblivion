package router

import "net/http"

type RequestHandler func(*Context) error

type handlerWrapper struct {
	handler RequestHandler
}

func (h *handlerWrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := &Context{
		Writer:  w,
		Request: req,
	}

	if err := h.handler(context); err != nil {
		writeErrorResponse(context, err)
		return
	}

	writeResponse(context)
}
