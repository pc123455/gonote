package framework

import (
	"context"
	"fmt"
	"gonote/framework/logger"
	"gonote/framework/route"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"
)

func stripHostPort(h string) string {
	// If no port on host, return unchanged
	if strings.IndexByte(h, ':') == -1 {
		return h
	}
	host, _, err := net.SplitHostPort(h)
	if err != nil {
		return h // on error, return unchanged
	}
	return host
}

// Return the canonical path for p, eliminating . and .. elements.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

type ErrorHandlerMap map[int]HandlerFunc

type Server struct {
	BindIP   string
	Port     int
	server   *http.Server
	handler  *Handler
	listener *Listener
	doneChan chan struct{}
	mu       sync.RWMutex
}

func (s *Server) Get(pattern string, handler HandlerFunc) {
	s.handler.AddHandleFunc("GET", pattern, handler)
}

func (s *Server) Post(pattern string, handler HandlerFunc) {
	s.handler.AddHandleFunc("POST", pattern, handler)
}

func (s *Server) Put(pattern string, handler HandlerFunc) {
	s.handler.AddHandleFunc("Put", pattern, handler)
}

func (s *Server) Delete(pattern string, handler HandlerFunc) {
	s.handler.AddHandleFunc("Delete", pattern, handler)
}

func (s *Server) Initialize(ip string, port int) error {
	s.doneChan = make(chan struct{})

	s.handler = &Handler{
		router: route.NewBaseRouter(),
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return err
	}
	s.listener = &Listener{
		ln: ln,
	}

	s.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        s.handler,
		MaxHeaderBytes: 1 << 30,
	}
	return nil
}

func (s *Server) Run() {
	err := s.server.ListenAndServe()
	logger.Errorf(err.Error())
	close(s.doneChan)
}

func (s *Server) Wait() {
	<-s.doneChan
}

func (s *Server) Shutdown() {
	s.server.Shutdown(context.Background())
	//close(s.doneChan)
}

type Listener struct {
	ln     net.Listener
	connWg sync.WaitGroup
}

func (l *Listener) Accept() (net.Conn, error) {
	l.connWg.Add(1)
	return l.ln.Accept()
}

func (l *Listener) Close() error {
	l.connWg.Done()
	return l.ln.Close()
}

func (l *Listener) Addr() net.Addr {
	return l.ln.Addr()
}

func (l *Listener) Wait() {
	l.connWg.Wait()
}
