package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cse224/tritonhttp"
)

func main() {
	currDir, err := os.Getwd() // Get current working directory, eg. /Users/username/go/src/cse224/tritonhttpd
	if err != nil {
		log.Fatalf("Could not get current working directory: %v", err)
	}
	// virtual_hosts.yaml: mapping from host name to the docRoot path
	default_vh_config_path := filepath.Join(currDir, "virtual_hosts.yaml")
	default_docroot := filepath.Join(currDir, "docroot_dirs")

	// Parse command line flags
	// eg. go run main.go -port=9090 -vh_config=virtual_hosts.yaml -docroot=docroot_dirs
	// no need to change vh_config and docroot_dirs_path
	var port = flag.Int("port", 8080, "the localhost port to listen on")
	var vh_config_path = flag.String("vh_config", default_vh_config_path, "path to the virtual hosting config file")
	var docroot_dirs_path = flag.String("docroot", default_docroot, "path to the directory that contains all docroot dirs")
	flag.Parse() // Parse command line flags, when called, it parses the command-line arguments from os.Args[1:]

	// Log server configs, print out the server configurations
	fmt.Println()
	log.Print("Server configs:")
	log.Printf("  port: %v", *port)
	log.Printf("  path to virtual hosts config file: %v", *vh_config_path)
	log.Printf("  path to docroot directories: %v", *docroot_dirs_path)
	fmt.Println()

	// Parse the virtual hosting config file, and return a map of host name to docRoot path
	// eg. virtual_hosts.yaml:
	// 	virtual_hosts:
	//		- hostName: "website1"
	//		docRoot: "htdocs1"
	// map[website1:/Users/username/go/src/cse224/tritonhttpd/docroot_dirs/htdocs1]
	virtualHosts := tritonhttp.ParseVHConfigFile(*vh_config_path, *docroot_dirs_path)

	// Start server
	// fmt.Sprintf: returns a formatted string, eg. ":9090"
	addr := fmt.Sprintf(":%v", *port)

	log.Printf("Starting TritonHTTP server")
	// server is listening on the port, and the virtualHosts map is passed to the server
	log.Printf("You can browse the website at http://localhost:%v/", *port)
	s := &tritonhttp.Server{
		Addr:         addr,
		VirtualHosts: virtualHosts,
	}
	// ListenAndServe listens on the TCP network address s.Addr and then handles requests on incoming connections
	log.Fatal(s.ListenAndServe())
}
