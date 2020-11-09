package server

import (
	"net/url"
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

// HandleFunc ...
type HandleFunc func(req *Request)

// Server ...
type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandleFunc
}

// Request ...
type Request struct {
	Conn net.Conn
	QueryParams url.Values
//	PathParams map[string]string
}
// NewServer ....
func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandleFunc)}
}

// Register .....
func (s *Server) Register(path string, handler HandleFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

// Start ....
func (s *Server) Start() error {

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		err = s.handle(&Request{Conn: conn, QueryParams: url.Values{}})
		if err != nil {
			log.Print(err)
			// Идём обслуживать следующего
			continue
		}

	}

}

func (s *Server) handle(req *Request) (err error) {
	defer func() {
		if cerr := req.Conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()
	// .....

	//conn.Write([]byte("Hello!\r\n"))

	buf := make([]byte, 4096)

	n, err := req.Conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
		return nil
	}

	if err != nil {
		log.Print(err)
		return err
	}
	//log.Printf("%s", buf[:n])

	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
	}

	method, path, version := parts[0], parts[1], parts[2]
	if method != "GET" {

	}

	if version != "HTTP/1.1" {

	}

	s.mu.RLock()
	for _, handler := range s.handlers {
		pth := path
		ind := bytes.Index([]byte(pth), []byte("?"))
		if ind != -1 {
			pth = pth[:ind]
		}

		params, err := url.ParseRequestURI(method + ":" + s.addr + ":/" + path)
		if err != nil {
			log.Print(err)
			break
		}

		req.QueryParams = url.Values(params.Query())

		handler(req)
		break
	}
	return nil

}
