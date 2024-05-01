package tritonhttp

import "time"

type Response struct {
	Proto      string // e.g. "HTTP/1.1"
	StatusCode int    // e.g. 200
	StatusText string // e.g. "OK"

	// Headers stores all headers to write to the response.
	Headers map[string]string

	// Request is the valid request that leads to this response.
	// It could be nil for responses not resulting from a valid request.
	// Hint: you might need this to handle the "Connection: Close" requirement
	Request *Request

	// FilePath is the local path to the file to serve.
	// It could be "", which means there is no file to serve.
	FilePath string
}

func (res *Response) HandleBadRequest() {
	res.Proto = "HTTP/1.1"
	res.StatusCode = 400
	res.StatusText = "Bad Request"
	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}
	res.Headers["Date"] = FormatTime(time.Now())
}

func (res *Response) HandleStatusNotFound() {
	res.Proto = "HTTP/1.1"
	res.StatusCode = 404
	res.StatusText = "Not Found"
	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}
	res.Headers["Date"] = FormatTime(time.Now())
}

func (res *Response) HandleOK() {
	res.Proto = "HTTP/1.1"
	res.StatusCode = 200
	res.StatusText = "OK"
	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}
	res.Headers["Date"] = FormatTime(time.Now())
}
