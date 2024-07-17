package httprouter

import (
	"net/http"
	"time"
)

type Context struct {
	Writer    http.ResponseWriter
	Request   *http.Request
	StartedAt time.Time
}

func (c *Context) fullRequestURI() string {
	return c.Request.Method + " " + c.Request.RequestURI
}

func (c *Context) duration() time.Duration {
	return time.Since(c.StartedAt)
}

// ***
// Context
// ***

func (c *Context) Deadline() (time.Time, bool) {
	return c.Request.Context().Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.Request.Context().Done()
}

func (c *Context) Err() error {
	return c.Request.Context().Err()
}

func (c *Context) Value(key any) any {
	return c.Request.Context().Value(key)
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
