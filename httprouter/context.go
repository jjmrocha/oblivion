package httprouter

import (
	"net/http"
	"time"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Start   time.Time
}

func (c *Context) fullRequestURI() string {
	return c.Request.Method + " " + c.Request.RequestURI
}

func (c *Context) duration() time.Duration {
	return time.Since(c.Start)
}

func (c *Context) OK(payload any) (*Response, error) {
	resp := Response{
		Status:  http.StatusOK,
		Payload: payload,
	}

	return &resp, nil
}

func (c *Context) Created(payload any) (*Response, error) {
	resp := Response{
		Status:  http.StatusCreated,
		Payload: payload,
	}

	return &resp, nil
}

func (c *Context) NoContent() (*Response, error) {
	resp := Response{
		Status:  http.StatusNoContent,
		Payload: nil,
	}

	return &resp, nil
}
