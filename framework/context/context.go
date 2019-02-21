package context

import (
	"net/http"
)

type HttpError struct {
	Status  int
	Message string
	Err     error
	ctx     *Context
}

func (this *HttpError) GetContext() *Context {
	return this.ctx
}

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
	isAbort bool
	Writer  http.ResponseWriter
	Error   HttpError
}

func (this *Response) WriteJson() {
	this.Writer.Header().Set("Content-Type", "application/json")
	if this.status == 0 {
		this.status = 200
	}
	//this.Writer.WriteHeader(this.status)
	if this.content != nil {
		this.Writer.Write(this.content)
	}
}

func (this *Response) Write() {
	this.Writer.Write(this.content)
}

func (this *Response) SetStatus(status int) {
	this.status = status
}

func (this *Response) GetStatus() int {
	return this.status
}

func (this *Response) AppendContent(b []byte) {
	this.content = append(this.content, b...)
}

func (this *Response) abort(err HttpError) {
	this.isAbort = true
	panic(err)
}

type Context struct {
	Input  Request
	Output Response
}

func (this *Context) Abort(err HttpError) {
	err.ctx = this
	this.Output.abort(err)
}
