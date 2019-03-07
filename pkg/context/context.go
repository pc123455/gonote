package context

import (
	"net/http"
)

type HttpError struct {
	Status  int
	Message []byte
	Err     error
	ctx     *Context
}

func (this *HttpError) GetContext() *Context {
	return this.ctx
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

func (this *Context) Abort(err HttpError) {
	err.ctx = this
	this.Output.abort(err)
}
