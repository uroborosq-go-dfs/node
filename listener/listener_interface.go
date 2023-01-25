package listener

import "github.com/uroborosq-go-dfs/node/service"

type NodeListener interface {
	Listen() error
}

func New(mode string, port string, serv *service.NodeService) NodeListener {
	if mode == "tcp" {
		return &TcpListener{
			port:    port,
			service: serv,
		}
	} else if mode == "http" {
		return &HttpListener{
			port:        port,
			nodeService: serv,
		}
	}
	return nil
}
