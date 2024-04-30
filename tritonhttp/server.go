package tritonhttp

import (
	"fmt"
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

	// Hint: create your listen socket and spawn off goroutines per incoming client
	// panic("todo")
	return nil
}
