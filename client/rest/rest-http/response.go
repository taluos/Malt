package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
	body []byte
}

func NewResponse(resp *http.Response) *Response {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return &Response{
		Response: resp,
		body:     body,
	}
}

func (r *Response) StatusCode() int {
	return r.Response.StatusCode
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) Header(key string) string {
	return r.Response.Header.Get(key)
}

func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal(r.body, v)
}

func (r *Response) String() string {
	return string(r.body)
}

func (r *Response) Reader() io.Reader {
	return bytes.NewReader(r.body)
}
