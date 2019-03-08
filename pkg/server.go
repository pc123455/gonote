package pkg

import (
	"context"
	"fmt"
	"gonote/pkg/logger"
	"gonote/pkg/route"
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
	BindIP  string
	Port    int
	server  *http.Server
	handler *Handler
	//listener *Listener
	doneChan chan struct{}
	mu       sync.RWMutex
}

func (s *Server) AppendStageHandlerFunc(stage int, handler HandlerFunc) {
	s.handler.AppendStageHandlerFunc(stage, handler)
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
	s.handler.Initialize()

	s.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        s.handler,
		MaxHeaderBytes: 1 << 30,
	}
	return nil
}

func (s *Server) GetDoneChan() <-chan struct{} {
	return s.doneChan
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
