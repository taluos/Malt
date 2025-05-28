package fasthttp

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/valyala/fasthttp"
)

type Response struct {
	*fasthttp.Response
}

func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal(r.Body(), v)
}

func (r *Response) String() string {
	return string(r.Body())
}

// 添加缺少的方法
func (r *Response) StatusCode() int {
	return r.Response.StatusCode()
}

func (r *Response) Body() []byte {
	return r.Response.Body()
}

func (r *Response) Header(key string) string {
	return string(r.Response.Header.Peek(key))
}

func (r *Response) Reader() io.Reader {
	return bytes.NewReader(r.Body())
}
