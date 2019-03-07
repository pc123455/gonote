package plugin

const (
	//protocol type
	protocolRPC = iota
	protocolHTTP
	protocolHTTPS

	//healthy type
	running = iota
	breakdown
)

type ProxyConfig struct {
}

type Backend struct {
	host     string
	port     int
	protocol int
	path     string
	health   int

	weight float32
}

type ProxyHandler struct {
	location string
	backEnds []Backend
}

type ProxyPlugin struct {
	handlers []ProxyHandler
}

func (pp *ProxyPlugin) Initialize(config string) {

}
