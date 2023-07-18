package server

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/config/static"
)

// EntrypointFactory the factory of TCP/UDP routers.
type EntryPointFactory struct {
	mu sync.Mutex

	staticConfiguration  static.Configuration
	staticEntryPointsTCP TCPEntryPoints
	staticEntryPointsUDP UDPEntryPoints
	routerFactory        *RouterFactory

	dynamicEntryPoints    static.EntryPoints
	dynamicEntryPointsTCP TCPEntryPoints
	dynamicEntryPointsUDP UDPEntryPoints
}

func NewEntryPointFactory(routerFactory *RouterFactory, config static.Configuration, tcpEntryPoints TCPEntryPoints, udpEntryPoints UDPEntryPoints) *EntryPointFactory {
	return &EntryPointFactory{
		routerFactory:         routerFactory,
		staticConfiguration:   config,
		staticEntryPointsTCP:  tcpEntryPoints,
		staticEntryPointsUDP:  udpEntryPoints,
		dynamicEntryPoints:    static.EntryPoints{},
		dynamicEntryPointsTCP: TCPEntryPoints{},
		dynamicEntryPointsUDP: UDPEntryPoints{},
	}
}

func (ef *EntryPointFactory) BuildEntryPoints(config dynamic.Configuration) {
	entryPoints := map[string]*static.EntryPoint{}

	ef.mu.Lock()
	defer ef.mu.Unlock()

	tcpEntryPointsForNewConf := map[string]*TCPEntryPoint{}
	for _, rt := range config.HTTP.Routers {
		for _, e := range rt.EntryPoints {
			if _, ok := ef.staticEntryPointsTCP[e]; ok {
				continue
			}
			if _, ok := ef.dynamicEntryPointsTCP[e]; ok {
				tcpEntryPointsForNewConf[e] = ef.dynamicEntryPointsTCP[e]
				continue
			}
			ep, ok := buildEntryPoint(e)
			if !ok {
				continue
			}

			tcpEntryPointsForNewConf[e] = ef.dynamicEntryPointsTCP[e]
			entryPoints[e] = ep
		}
	}
	for _, rt := range config.TCP.Routers {
		for _, e := range rt.EntryPoints {
			if _, ok := ef.staticEntryPointsTCP[e]; ok {
				continue
			}
			if _, ok := ef.dynamicEntryPointsTCP[e]; ok {
				tcpEntryPointsForNewConf[e] = ef.dynamicEntryPointsTCP[e]
				continue
			}
			ep, ok := buildEntryPoint(e)
			if !ok {
				continue
			}
			tcpEntryPointsForNewConf[e] = ef.dynamicEntryPointsTCP[e]
			entryPoints[e] = ep
		}
	}

	udpEntryPointsForNewConf := map[string]*UDPEntryPoint{}
	for _, rt := range config.UDP.Routers {
		for _, e := range rt.EntryPoints {
			if _, ok := ef.staticEntryPointsUDP[e]; ok {
				continue
			}
			if _, ok := ef.dynamicEntryPointsUDP[e]; ok {
				udpEntryPointsForNewConf[e] = ef.dynamicEntryPointsUDP[e]
				continue
			}
			ep, ok := buildEntryPoint(e)
			if !ok {
				continue
			}
			udpEntryPointsForNewConf[e] = ef.dynamicEntryPointsUDP[e]
			entryPoints[e] = ep
		}
	}

	deletedEntryPointsTCP := TCPEntryPoints{}
	for name, ep := range ef.dynamicEntryPointsTCP {
		if _, ok := tcpEntryPointsForNewConf[name]; !ok {
			deletedEntryPointsTCP[name] = ep
			delete(ef.dynamicEntryPointsTCP, name)
			delete(ef.dynamicEntryPoints, name)
		}
	}

	deletedEntryPointsUDP := UDPEntryPoints{}
	for name, ep := range ef.dynamicEntryPointsUDP {
		if _, ok := udpEntryPointsForNewConf[name]; !ok {
			deletedEntryPointsUDP[name] = ep
			delete(ef.dynamicEntryPointsUDP, name)
			delete(ef.dynamicEntryPoints, name)
		}
	}

	for _, e := range entryPoints {
		e.SetDefaults()
	}

	if len(entryPoints) > 0 {
		for n, e := range entryPoints {
			ef.dynamicEntryPoints[n] = e
		}
	}

	newEntryPointsTCP := NewTCPEntryPointsIgnoreErr(entryPoints, ef.staticConfiguration.HostResolver)
	if len(newEntryPointsTCP) > 0 {
		newEntryPointsTCP.Start()
		for n, e := range newEntryPointsTCP {
			ef.dynamicEntryPointsTCP[n] = e
		}
	}

	newEntryPointsUDP := NewUDPEntryPointsIgnoreErr(entryPoints)
	if len(newEntryPointsUDP) > 0 {
		newEntryPointsUDP.Start()
		for n, e := range newEntryPointsUDP {
			ef.dynamicEntryPointsUDP[n] = e
		}
	}

	if len(deletedEntryPointsTCP) > 0 {
		deletedEntryPointsTCP.Stop()
	}

	if len(deletedEntryPointsUDP) > 0 {
		deletedEntryPointsUDP.Stop()
	}

	ef.updateRouterFactory()
}

func (ef *EntryPointFactory) ServerEntryPointsTCP() TCPEntryPoints {
	eps := make(TCPEntryPoints, len(ef.staticEntryPointsTCP)+len(ef.dynamicEntryPointsTCP))
	for key, ep := range ef.staticEntryPointsTCP {
		eps[key] = ep
	}
	for key, ep := range ef.dynamicEntryPointsTCP {
		eps[key] = ep
	}
	return eps

}

func (ef *EntryPointFactory) ServerEntryPointsUDP() UDPEntryPoints {
	eps := make(UDPEntryPoints, len(ef.staticEntryPointsUDP)+len(ef.dynamicEntryPointsUDP))
	for key, ep := range ef.staticEntryPointsUDP {
		eps[key] = ep
	}
	for key, ep := range ef.dynamicEntryPointsUDP {
		eps[key] = ep
	}
	return eps
}

func (ef *EntryPointFactory) updateRouterFactory() {
	eps := static.EntryPoints{}
	for n, e := range ef.dynamicEntryPoints {
		eps[n] = e
	}
	for n, e := range ef.staticConfiguration.EntryPoints {
		eps[n] = e
	}

	ef.routerFactory.UpdateEntryPoints(eps)
}

func buildEntryPoint(ep string) (e *static.EntryPoint, ok bool) {
	strs := strings.Split(ep, "-")
	if len(strs) != 2 {
		return nil, false
	}

	if !(checkPort(strs[1]) && checkProtocol(strs[0])) {
		return nil, false
	}

	return &static.EntryPoint{
		Address: fmt.Sprintf(":%s/%s", strs[1], toProtocol(strs[0])),
	}, true
}

const (
	tcpProto  = "tcp"
	udpProto  = "udp"
	httpProto = "http"
)

func checkProtocol(p string) bool {
	return p == httpProto || p == tcpProto || p == udpProto
}

func checkPort(p string) bool {
	if port, err := strconv.Atoi(p); err == nil {
		return port > 0 && port <= 65535
	}
	return false
}

func toProtocol(p string) string {
	switch p {
	case httpProto:
		return tcpProto
	case udpProto:
		return udpProto
	default:
		return tcpProto
	}
}
