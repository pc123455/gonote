package context

import (
	"net/http"
)

type Request struct {
	//Header http.Header
	//Method string
	//URL *url.URL
	//Data []byte
	//Body io.ReadCloser
	*http.Request
	Args       map[string]interface{}
	RawContent []byte
}

func (this *Request) ReadBody() []byte {
	var content []byte
	this.Body.Read(content)
	this.RawContent = content
	return content
}

type Response struct {
	//Header http.Header
	//File http.File
	content []byte
	status  int
	Writer  http.ResponseWriter
}

func (this *Response) ServeJson() {
	this.Writer.Header().Set("Content-Type", "application/json")
	if this.status == 0 {
		this.status = 200
	}
	//this.Writer.WriteHeader(this.status)
	if this.content != nil {
		this.Writer.Write(this.content)
	}
}

func (this *Response) SetStatus(status int) {
	this.status = status
}

func (this *Response) Write(b []byte) {
	this.content = append(this.content, b...)
}

type Context struct {
	Input  Request
	Output Response
}
