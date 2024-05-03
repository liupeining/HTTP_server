package tritonhttp

import (
	"bufio"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

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
		res.Headers["Close"] = "close"
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
	res.Headers["Close"] = "close"
	res.Headers["Date"] = FormatTime(time.Now())
}

func (res *Response) HandleOK(docRoot string, req *Request) {
	res.Request = req
	res.Proto = "HTTP/1.1"
	res.StatusCode = 200
	res.StatusText = "OK"
	if res.Headers == nil {
		res.Headers = make(map[string]string)
	}
	res.Headers["Date"] = FormatTime(time.Now())
	res.FilePath = docRoot + res.Request.URL
	// If URL ends with /, interpret as index.html
	if res.Request.URL[len(res.Request.URL)-1] == '/' {
		res.FilePath += "index.html"
	}
	fmt.Println("File Path: ", res.FilePath)

	// Check if the file exists
	if _, err := os.Stat(res.FilePath); os.IsNotExist(err) {
		fmt.Println("File does not exist")
		res.HandleStatusNotFound()
		//req.Headers["Close"] = "close"
		return
	}

	stats, err := os.Stat(res.FilePath)
	if err != nil {
		fmt.Println("Error in getting file stats: ", err)
		res.HandleStatusNotFound()
		//req.Headers["Close"] = "close"
		return
	}

	res.Headers["Content-Length"] = strconv.FormatInt(stats.Size(), 10)
	res.Headers["Content-Type"] = mime.TypeByExtension(filepath.Ext(res.FilePath))
	res.Headers["Date"] = FormatTime(time.Now())
	res.Headers["Last-Modified"] = FormatTime(stats.ModTime())
}

func (res *Response) Write(w io.Writer) error {
	// Write the response line
	bw := bufio.NewWriter(w)
	statusLine := fmt.Sprintf("%v %v %v\r\n", res.Proto, res.StatusCode, res.StatusText)
	if _, err := bw.WriteString(statusLine); err != nil {
		return err
	}
	// Write Headers
	headers := res.Headers
	headerKeys := make([]string, 0)
	for key, _ := range headers {
		headerKeys = append(headerKeys, key)
	}
	sort.Strings(headerKeys)

	for _, key := range headerKeys {
		keyValue := key + ": " + headers[key] + "\r\n"
		if _, err := bw.WriteString(keyValue); err != nil {
			return err
		}
	}
	if _, err := bw.WriteString("\r\n"); err != nil {
		return err
	}

	filePath := res.FilePath
	if len(filePath) > 0 {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		if _, err := bw.Write(data); err != nil {
			return err
		}
	}
	if err := bw.Flush(); err != nil {
		// log.Println("Flush error: ", err)
		return nil
	}
	return nil
}
