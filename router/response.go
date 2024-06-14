package router

import "net/http"

type Response struct {
	Status  int
	Payload any
}

type RequestHandler func(*Context) (*Response, error)

func (h RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := &Context{
		Writer:  w,
		Request: req,
	}

	resp, err := h(context)
	if err != nil {
		writeErrorResponse(context, err)
		return
	}

	writeResponse(context, resp)
}
