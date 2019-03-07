package context

import "net/http"

type Response struct {
	//Header http.Header
	//File http.File
	content []byte
	status  int
	isAbort bool
	Writer  http.ResponseWriter
	Error   *HttpError
}

func newOutput(writer http.ResponseWriter) Response {
	return Response{
		content: []byte{},
		status:  200,
		isAbort: false,
		Writer:  writer,
		Error:   nil,
	}
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

func (this *Response) Write() {
	this.Writer.WriteHeader(this.status)
	this.Writer.Write(this.content)
}

func (this *Response) Flush() {
	this.status = 200
	this.content = make([]byte, 0)
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
	this.Error = &err
	panic(err)
}
