package tritonhttp

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Server struct {
	// Addr specifies the TCP address for the server to listen on,
	// in the form "host:port". It shall be passed to net.Listen()
	// during ListenAndServe().
	Addr string // e.g. ":0"

	// VirtualHosts contains a mapping from host name to the docRoot path
	// (i.e. the path to the directory to serve static files from) for
	// all virtual hosts that this server supports
	VirtualHosts map[string]string
}

// ValidateServerSetup checks the validity of the docRoot of the server
func (s *Server) ValidateServerSetup() error {
	// Validating the doc root of the server
	for _, docRoot := range s.VirtualHosts {
		fi, err := os.Stat(docRoot)
		if os.IsNotExist(err) {
			return err
		}
		if !fi.IsDir() {
			return fmt.Errorf("doc root %q is not a directory", docRoot)
		}
	}
	return nil
}

// ListenAndServe listens on the TCP network address s.Addr and then
// handles requests on incoming connections.
func (s *Server) ListenAndServe() error {
	// Hint: Validate all docRoots
	if err := s.ValidateServerSetup(); err != nil {
		return fmt.Errorf("server is not setup correctly %v", err)
	}
	addr := "localhost" + s.Addr
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error in listening on : %v", addr, err)
	}
	fmt.Println("Listening on", ln.Addr())
	defer func() {
		err = ln.Close()
		if err != nil {
			fmt.Println("error in closing listener", err)
		}
	}()
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go s.HandleConnection(conn)
	}
}

// HandleConnection reads requests from the accepted conn and handles them.
func (s *Server) HandleConnection(conn net.Conn) {
	br := bufio.NewReader(conn)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			log.Printf("Failed to set timeout for connection %v", conn)
			_ = conn.Close()
			return
		}
		// Read request from the client
		var isEOF bool
		req, err, isEOF := ReadRequest(br)
		if isEOF {
			return
		}
		fmt.Println("Request: ", req)

		res := &Response{}
		res.Headers = make(map[string]string)
		if req != nil && req.Close {
			res.Headers["Connection"] = "close"
		}
		if err != nil {
			res.HandleBadRequest()
			err = res.Write(conn)
			if err != nil {
				fmt.Println("Error in writing response(400): ", err)
			}
			fmt.Println("Response: ", res)
			_ = conn.Close()
			return
		}
		res.HandleOK(s.VirtualHosts[req.Host], req) // pass the docRoot of the host to HandleOK
		err = res.Write(conn)
		if err != nil {
			fmt.Println("Error in writing response: ", err)
		}
		fmt.Println("Response: ", res)

		if res.Headers["Connection"] == "close" {
			_ = conn.Close()
			return
		}
	}
}

// ReadRequest reads and parses a request from the buffered reader.
func ReadRequest(br *bufio.Reader) (req *Request, err error, isEOF bool) {
	req = &Request{} // Method, URL, Proto, Headers, Host, Close
	// Read the first line of the request, which contains the method, URL, and protocol eg. GET /index.html HTTP/1.1
	firstLine, err := br.ReadString('\n')
	if err != nil {
		if err.Error() == "EOF" {
			return nil, err, true
		}
		fmt.Println("Error in reading first line: ", err)
		return nil, err, false
	}
	err = parseFirstLine(firstLine, req)
	if err != nil {
		fmt.Println("Error in parsing first line: ", err)
		return nil, err, false
	}
	req.Headers = make(map[string]string)
	err = parseHeaders(br, req)
	if err != nil {
		fmt.Println("Error in parsing headers: ", err)
		return nil, err, false
	}
	req.Host = req.Headers["Host"]
	req.Close = req.Headers["Connection"] == "close"
	return req, nil, false
}

func parseHeaders(br *bufio.Reader, req *Request) error {
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			fmt.Println("Error in reading line: ", err)
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid header: %q", line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		req.Headers[key] = value
	}
	return nil
}

func parseFirstLine(firstLine string, req *Request) error {
	parts := strings.Split(firstLine, " ")
	if len(parts) != 3 {
		return fmt.Errorf("invalid first line: %q", firstLine)
	}
	if parts[0] != "GET" {
		return fmt.Errorf("invalid method: %q", parts[0])
	}
	req.Method = parts[0]
	if !strings.HasPrefix(parts[1], "/") {
		return fmt.Errorf("invalid URL: %q", parts[1])
	}
	req.URL = parts[1]
	protocol := strings.TrimSpace(parts[2])
	if protocol != "HTTP/1.1" {
		return fmt.Errorf("invalid protocol: %q", parts[2])
	}
	req.Proto = protocol
	return nil
}
