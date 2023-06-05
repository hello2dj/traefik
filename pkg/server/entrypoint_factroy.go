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
		routerFactory:        routerFactory,
		staticConfiguration:  config,
		staticEntryPointsTCP: tcpEntryPoints,
		staticEntryPointsUDP: udpEntryPoints,
	}
}

func (ef *EntryPointFactory) BuildEntryPoints(config dynamic.Configuration) {
	entryPoints := map[string]*static.EntryPoint{}

	ef.mu.Lock()
	defer ef.mu.Unlock()

	for _, rt := range config.HTTP.Routers {
		for _, e := range rt.EntryPoints {
			if _, ok := ef.staticEntryPointsTCP[e]; ok {
				continue
			}
			if _, ok := ef.dynamicEntryPointsTCP[e]; ok {
				continue
			}
			ep, ok := buildEntryPoint(e)
			if !ok {
				continue
			}
			entryPoints[e] = ep
		}
	}
	for _, rt := range config.TCP.Routers {
		for _, e := range rt.EntryPoints {
			if _, ok := ef.staticEntryPointsTCP[e]; ok {
				continue
			}
			if _, ok := ef.dynamicEntryPointsTCP[e]; ok {
				continue
			}
			ep, ok := buildEntryPoint(e)
			if !ok {
				continue
			}
			entryPoints[e] = ep
		}
	}

	for _, rt := range config.UDP.Routers {
		for _, e := range rt.EntryPoints {
			if _, ok := ef.staticEntryPointsUDP[e]; ok {
				continue
			}
			if _, ok := ef.dynamicEntryPointsUDP[e]; ok {
				continue
			}
			ep, ok := buildEntryPoint(e)
			if !ok {
				continue
			}
			entryPoints[e] = ep
		}
	}

	deletedEntryPointsTCP := TCPEntryPoints{}
	for name, ep := range ef.dynamicEntryPointsTCP {
		if _, ok := entryPoints[name]; !ok {
			deletedEntryPointsTCP[name] = ep
			delete(ef.dynamicEntryPointsTCP, name)
		}
	}

	deletedEntryPointsUDP := UDPEntryPoints{}
	for name, ep := range ef.dynamicEntryPointsUDP {
		if _, ok := entryPoints[name]; !ok {
			deletedEntryPointsUDP[name] = ep
			delete(ef.dynamicEntryPointsUDP, name)
		}
	}

	for _, e := range entryPoints {
		e.SetDefaults()
	}

	ef.dynamicEntryPoints = entryPoints
	newEntryPointsTCP := NewTCPEntryPointsIgnoreErr(entryPoints, ef.staticConfiguration.HostResolver)
	ef.dynamicEntryPointsTCP = newEntryPointsTCP

	newEntryPointsUDP := NewUDPEntryPointsIgnoreErr(entryPoints)
	ef.dynamicEntryPointsUDP = newEntryPointsUDP

	deletedEntryPointsTCP.Stop()
	deletedEntryPointsUDP.Stop()

	newEntryPointsTCP.Start()
	newEntryPointsUDP.Start()

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
