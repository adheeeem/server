package server

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
)

var ErrVersionNotCompatible = errors.New("not compatible version")

type HandlerFunc func(req *Request)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

type Request struct {
	Conn        net.Conn
	QueryParams url.Values
	PathParams  map[string]string
	Headers     map[string]string
	Body        []byte
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}
func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(cerr)
		}
	}()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Print(err)
			continue
		}

		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) (err error) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			if err == nil {
				err = cerr
				return
			}
			log.Print(err)
		}
	}()
	var reqPath string
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)

	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	data := buf[:n]

	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)

	if requestLineEnd == -1 {
		return
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return
	}

	var req Request

	method, path, version := parts[0], parts[1], parts[2]
	uri, err := url.ParseRequestURI(path)

	id := uri.Query().Get("id")

	if method == "POST" || method == "PUT" {
		bodyDelim := []byte{'{'}
		brStart := bytes.Index(data, bodyDelim)
		body := data[brStart+1 : len(data)-2]
		req.Body = body
		log.Print(string(req.Body))
	}

	if id != "" {
		req.QueryParams = map[string][]string{
			"id": {id},
		}
		reqPath = "/payments"
	} else {
		reqPars := strings.Split(uri.String(), "/")
		req.PathParams = map[string]string{
			"id": reqPars[2],
		}
		reqPath = "/payments/{id}"
	}

	req.Conn = conn

	if version != "HTTP/1.1" {
		return ErrVersionNotCompatible
	}

	s.mu.RLock()

	handlePath, ok := s.handlers[reqPath]

	s.mu.RUnlock()

	if ok {
		req.Headers = getHeaders(string(data))
		handlePath(&req)
	}

	return nil
}

func getHeaders(header string) map[string]string {
	lines := strings.Split(header, "\n")
	result := make(map[string]string)
	for _, line := range lines {
		temp := strings.Split(line, ":")
		if len(temp) == 2 {
			result[temp[0]] = temp[1]
		}
	}
	return result
}
