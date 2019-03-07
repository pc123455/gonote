package context

import "net/http"

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
