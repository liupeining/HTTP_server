package tritonhttp

import (
	"fmt"
	"net"
	"os"
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
	// fi, err := os.Stat(s.DocRoot)
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
	fmt.Println("Server setup valid!")

	// server should now start to listen on the configured address
	addr := "localhost" + s.Addr
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error in listening on : %v", addr, err)
	}
	fmt.Println("Listening on", ln.Addr())

	// making sure the listener is closed when we exit
	defer func() {
		err = ln.Close()
		if err != nil {
			fmt.Println("error in closing listener", err)
		}
	}()

	// accept connections forever
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		fmt.Println("accepted connection", conn.RemoteAddr())
		go s.HandleConnection(conn)
		if conn != nil {
			conn.Close()
			break
		}
	}
	return nil
}

// HandleConnection reads requests from the accepted conn and handles them.
func (s *Server) HandleConnection(conn net.Conn) {
	// now doing nothing
	_ = conn.Close()
}
