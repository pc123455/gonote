package context

import (
	"net/http"
)

type Context struct {
	http.ResponseWriter
	*http.Request
	Param *Param
}
