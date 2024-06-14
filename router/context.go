package router

import "net/http"

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
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
