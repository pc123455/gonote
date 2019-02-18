package route

import (
	"gonote/framework"
	"gonote/framework/logger"
	"regexp"
	"strings"
)

type pathNode struct {
	handler func(ctx *framework.Context)
	childs  map[string]*pathNode
}

func (this *pathNode) initialize(handler func(ctx *framework.Context)) {
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

func (this *pathNode) getHandler() func(ctx *framework.Context) {
	return this.handler
}

func (this *pathNode) setHandler(handler func(ctx *framework.Context)) {
	this.handler = handler
}

func (this *pathNode) match(pathSequence []string) (handler func(ctx *framework.Context)) {
	handler = nil
	pathWord := pathSequence[0]
	child := this.childs[pathWord]
	variableReg := regexp.MustCompile("^<[\\w:]*>$")

	if child != nil {
		if len(pathSequence) > 1 {
			//if current path node is not leaf node
			handler = child.match(pathSequence[1:])
		} else {
			if child.getHandler() == nil {
				return handler
			}
		}
	}

	if len(pathSequence) == 1 {
		return handler
	}
	for k, v := this.childs {
		if variableReg.Match(k) {

		}
	}
}

type VariableRouter struct {
	variableRouteMethodMap map[string]*pathNode
}

func (this *VariableRouter) match(node *pathNode, pathSequence []string) {

}

func (this *VariableRouter) AddRoute(method string, pattern string, handler func(ctx *framework.Context)) error {
	_, ok := this.variableRouteMethodMap[method]
	var rootNode *pathNode
	if !ok {
		rootNode = new(pathNode)
		rootNode.initialize(nil)
		this.variableRouteMethodMap[method] = rootNode
	}

	pathSequence := strings.Split(pattern, "/")
	currentNode := rootNode
	for _, pathWord := range pathSequence {
		//iterate node until leaf path
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
		logger.Warnf("method: %s and path: %s already exist", method, pattern)
	}
	return nil
}

func (this *VariableRouter) MatchRoute(method string, path string) (handler func(ctx *framework.Context), param Param) {
	rootNode, ok := this.variableRouteMethodMap[method]
	if rootNode
	return nil
}

func (this *VariableRouter) Initialize() {

}
