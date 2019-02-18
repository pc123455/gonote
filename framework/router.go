package framework

import (
	"fmt"
	"net/http"
	"strings"
)

type router interface {
	addRoute(method string, pattern string, handler func(ctx *Context)) error
	getRoute(method string, path string) func(ctx *Context)
}

type baseRouter struct {
	methodMap map[string]*pathNode
	//pathHandleMap map[string]func(ctx *context.Context)
}

type pathNode struct {
	handler func(ctx *Context)
	childs  map[string]*pathNode
}

func (this *pathNode) initialize(handler func(ctx *Context)) {
	this.handler = handler
	this.childs = make(map[string]*pathNode, 10)
}

func (this *pathNode) getChild(key string) *pathNode {
	child, ok := this.childs[key]
	if ok {
		return child
	}
	return nil
}

func (this *pathNode) setChild(key string, node *pathNode) {
	this.childs[key] = node
}

func (this *pathNode) getHandler() func(ctx *Context) {
	return this.handler
}

func (this *pathNode) setHandler(handler func(ctx *Context)) {
	this.handler = handler
}

func newBaseRouter() baseRouter {
	return baseRouter{make(map[string]*pathNode, 6)}
}

func (this baseRouter) addRoute(method string, pattern string, handler func(ctx *Context)) error {
	//regex := "(/[a-zA-Z0-9._~-])*(/<[a-zA-Z]+>)(/[a-zA-Z0-9._~-])*"
	method = strings.ToUpper(strings.TrimSpace(method))
	_, ok := this.methodMap[method]
	var rootNode *pathNode
	if !ok {
		rootNode = new(pathNode)
		rootNode.initialize(nil)
		this.methodMap[method] = rootNode
	}

	pathSequence := strings.Split(pattern, "/")
	currentNode := rootNode
	for _, pathWord := range pathSequence {
		nextNode := currentNode.getChild(pathWord)
		if nextNode == nil {
			nextNode := new(pathNode)
			nextNode.initialize(nil)
			currentNode.setChild(pathWord, nextNode)
		}
		currentNode = nextNode
	}

	if currentNode.getHandler() == nil {
		currentNode.setHandler(handler)
	} else {

	}

	pathMap, _ := this.methodMap[method]

	pattern = strings.TrimSpace(pattern)
	_, ok = pathMap[pattern]
	if ok {
		panic(fmt.Sprintf("method: %s and path: %s already exist", method, pattern))
	}

	pathMap[pattern] = handler

	return nil
}

func (this baseRouter) getRoute(method string, path string) func(ctx *Context) {
	method = strings.ToUpper(strings.TrimSpace(method))
	pathMap, ok := this.methodMap[method]
	if !ok {
		return nil
	}

	handler, ok := pathMap[path]
	if ok {
		return handler
	}

	return nil
}
