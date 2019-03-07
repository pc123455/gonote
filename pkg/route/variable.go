package route

import (
	"gonote/pkg/context"
	"gonote/pkg/logger"
	"regexp"
	"strconv"
	"strings"
)

type pathNode struct {
	handler func(ctx *context.Context)
	childs  map[string]*pathNode
}

func (this *pathNode) initialize(handler func(ctx *context.Context)) {
	this.handler = handler
	this.childs = nil
}

func (this *pathNode) getChild(key string) *pathNode {
	if this.childs == nil {
		return nil
	}
	return this.childs[key]
}

func (this *pathNode) setChild(key string, node *pathNode) {
	if this.childs == nil {
		this.childs = map[string]*pathNode{}
	}
	this.childs[key] = node
}

func (this *pathNode) getHandler() func(ctx *context.Context) {
	return this.handler
}

func (this *pathNode) setHandler(handler func(ctx *context.Context)) {
	this.handler = handler
}

func (this *pathNode) match(pathSequence []string) (handler func(ctx *context.Context), param context.Param) {
	handler = nil
	param = make(context.Param)
	if len(pathSequence) <= 0 {
		// if current path node is the end of path, return handler of current node
		return this.handler, param
	}

	if this.childs == nil {
		// current node is a leaf node
		return nil, param
	} else {
		pathWord := pathSequence[0]
		child := this.childs[pathWord]
		variableReg := regexp.MustCompile("^<[\\w:-]*>$")

		if child != nil {
			if len(pathSequence) > 0 {
				var childParam context.Param
				handler, childParam = child.match(pathSequence[1:])
				if handler != nil {
					context.MergeParam(param, childParam)
				}
			}
		}

		//match variable nodes
		for k, v := range this.childs {
			if handler != nil {
				return handler, param
			}
			if variableReg.Match([]byte(k)) {
				if len(pathWord) > 0 {
					var childParam context.Param
					key := k[1 : len(k)-1]
					if strings.Contains(key, ":") {
						keySeq := strings.Split(key, ":")
						paramName := keySeq[1]
						switch paramType := keySeq[0]; paramType {
						case "int":
							num, err := strconv.Atoi(pathWord)
							if err == nil {
								param[paramName] = num
							}
						case "float":
							num, err := strconv.ParseFloat(pathWord, 64)
							if err == nil {
								param[paramName] = num
							}
						default:
							param[paramName] = pathWord
						}
					} else {
						param[key] = pathWord
					}

					handler, childParam = v.match(pathSequence[1:])
					if handler != nil {
						context.MergeParam(param, childParam)
					}
				}
			}
		}
	}
	return handler, param

}

type VariableRouter struct {
	variableRouteMethodMap map[string]*pathNode
}

func (this *VariableRouter) match(node *pathNode, pathSequence []string) {

}

func (this *VariableRouter) AddRoute(method string, pattern string, handler func(ctx *context.Context)) error {
	rootNode := this.variableRouteMethodMap[method]
	if rootNode == nil {
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
			nextNode = new(pathNode)
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

func (this *VariableRouter) MatchRoute(method string, path string) (handler func(ctx *context.Context), param context.Param) {
	rootNode := this.variableRouteMethodMap[method]
	param = context.Param{}
	if rootNode == nil {
		return nil, param
	}
	path = strings.TrimSpace(path)
	pathSequence := strings.Split(path, "/")
	return rootNode.match(pathSequence)
}

func (this *VariableRouter) Initialize() {
	this.variableRouteMethodMap = make(map[string]*pathNode, 6)
}
