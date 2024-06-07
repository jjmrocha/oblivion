package router

import "net/http"

type response struct {
	status  int
	payload any
}

type RequestHandler func(*Context) error

type handlerWrapper struct {
	handler RequestHandler
}

func (h *handlerWrapper) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := Context{
		Writer:  w,
		Request: req,
	}

	if err := h.handler(&context); err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeResponse(w, context.response.status, context.response.payload)
}
