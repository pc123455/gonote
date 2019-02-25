package framework

import (
	"fmt"
	"gonote/framework/route"
	"net"
	"net/http"
	"path"
	"strings"
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
}

func (this *Server) Get(pattern string, handler HandlerFunc) {
	this.handler.AddHandleFunc("GET", pattern, handler)
}

func (this *Server) Post(pattern string, handler HandlerFunc) {
	this.handler.AddHandleFunc("POST", pattern, handler)
}

func (this *Server) Put(pattern string, handler HandlerFunc) {
	this.handler.AddHandleFunc("Put", pattern, handler)
}

func (this *Server) Delete(pattern string, handler HandlerFunc) {
	this.handler.AddHandleFunc("Delete", pattern, handler)
}

func (this *Server) Initialize(ip string, port int) {
	//this.router = route.NewBaseRouter()

	this.handler = &Handler{
		router: route.NewBaseRouter(),
	}

	this.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		Handler:        this.handler,
		MaxHeaderBytes: 1 << 30,
	}
}

func (this *Server) Run() {
	this.server.ListenAndServe()
}
