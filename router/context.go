package router

import "net/http"

type response struct {
	status  int
	payload any
}

type Context struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	response *response
}

func (c *Context) OK(payload any) error {
	response := response{
		status:  http.StatusOK,
		payload: payload,
	}
	c.response = &response
	return nil
}

func (c *Context) Created(payload any) error {
	response := response{
		status:  http.StatusCreated,
		payload: payload,
	}
	c.response = &response
	return nil
}

func (c *Context) NoContent() error {
	response := response{
		status:  http.StatusNoContent,
		payload: nil,
	}
	c.response = &response
	return nil
}
