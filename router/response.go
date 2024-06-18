package router

import (
	"log"
	"net/http"
	"time"
)

type Response struct {
	Status  int
	Payload any
}

type RequestHandler func(*Context) (*Response, error)

func (h RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Writer:  w,
		Request: req,
		Start:   time.Now(),
	}

	resp, err := h(ctx)
	if err != nil {
		log.Printf("ERROR => %s => %v", ctx.fullRequestURI(), err.Error())
		resp = errorResponse(err)
	}

	writeResponse(ctx, resp)
}
