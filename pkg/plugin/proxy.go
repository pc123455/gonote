package plugin

import (
	"encoding/json"
	"fmt"
	"gonote/pkg"
	"gonote/pkg/context"
	"gonote/pkg/logger"
	"net/url"
)

const (
	//protocol type
	protocolRPC = iota
	protocolHTTP
	protocolHTTPS

	//healthy type
	running = iota
	breakdown

	//load balance type
	rotationForward = iota
	ipHash
)

var (
	config ProxyConfig
)

type ProxyConfig struct {
	Rules []ProxyRule
}

type ProxyRule struct {
	Pattern  string
	IPHash   bool
	Backends []Backend
}

func (pr *ProxyRule) init() error {
	length := len(pr.Backends)

	for i := 0; i < length; i++ {
		err := pr.Backends[i].parseUri()
		if err != nil {
			return err
		}
	}
	return nil
}

func (pr *ProxyRule) match(url *url.URL) bool {
	return pr.prefixMatch(url)
}

func (pr *ProxyRule) prefixMatch(url *url.URL) bool {
	//prefix match
	if len(pr.Pattern) > len(url.Path) {
		return false
	}
	if pr.Pattern == url.Path[:len(pr.Pattern)] {
		return true
	}
	return false
}

func (pr *ProxyRule) upRequest(ctx *context.Context) error {
	return nil
}

type Backend struct {
	url    *url.URL
	health int

	Uri    string
	Weight float32
}

func (b *Backend) parseUri() (err error) {
	if b.url, err = url.Parse(b.Uri); err != nil {
		return fmt.Errorf("uri parsing error: %s, url is %s", err, b.Uri)
	}
	return
}

type ProxyHandler struct {
	pattern  string
	lbType   int
	backEnds []Backend
}

func Init(jsonConf string, server *pkg.Server) {
	// parse config
	err := json.Unmarshal([]byte(jsonConf), config)
	if err != nil {
		logger.Errorf("config unmarshal error: %s", err)
	}

	rules := config.Rules
	ruleLen := len(rules)

	for i := 0; i < ruleLen; i++ {
		rules[i].init()
	}

	//register handle
	server.AppendStageHandlerFunc(pkg.BeforeRouteStage, Handler)
}

func Handler(ctx *context.Context) {
	var matchedRule *ProxyRule
	for _, rule := range config.Rules {
		if rule.match(ctx.Input.URL) {
			matchedRule = &rule
		}
	}
	matchedRule.upRequest(ctx)
}
