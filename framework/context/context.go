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
	Args       Param
	RawContent []byte
}

func newInput(r *http.Request) Request {
	input := Request{
		Request:    r,
		Args:       make(Param),
		RawContent: nil,
	}
	return input
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
	Input     Request
	Output    Response
	stageOver bool
	//stage int
	//handlerIndex int
}

func NewContext(writer http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Input:  newInput(r),
		Output: newOutput(writer),
	}
}

func (this *Context) NextStage() {
	this.stageOver = true
}

func (this *Context) IsStageOver() bool {
	return this.stageOver
}

func (this *Context) ResetStageOver() {
	this.stageOver = false
}

//
//func (this *Context) GetStage() int {
//	return this.stage
//}
//
//func (this *Context) GetHandlerIndex() int {
//	return this.handlerIndex
//}
//
//func (this *Context) IncreaseHandlerIndex() {
//	this.handlerIndex++
//}

func (this *Context) Abort(err HttpError) {
	err.ctx = this
	this.Output.abort(err)
}
