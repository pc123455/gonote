package framework

import "net/http"

type Context struct {
	http.ResponseWriter
	*http.Request
}
