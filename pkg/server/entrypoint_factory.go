package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/log"
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

func NewEntryPointFactory(
	routerFactory *RouterFactory,
	config static.Configuration,
	tcpEntryPoints TCPEntryPoints,
	udpEntryPoints UDPEntryPoints,
) *EntryPointFactory {
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

	for _, rt := range config.HTTP.Routers {
		for _, e := range rt.EntryPoints {
			if _, ok := ef.staticEntryPointsTCP[e]; ok {
				assignEntryPointTransport(rt, ef.staticEntryPointsTCP[e])
			}
			if _, ok := ef.dynamicEntryPointsTCP[e]; ok {
				assignEntryPointTransport(rt, ef.dynamicEntryPointsTCP[e])
			}

			if _, ok := ef.dynamicEntryPoints[e]; ok {
				assignEntryPointTransportForStatic(rt, ef.dynamicEntryPoints[e])
			}
			if _, ok := ef.staticConfiguration.EntryPoints[e]; ok {
				assignEntryPointTransportForStatic(rt, ef.staticConfiguration.EntryPoints[e])
			}
		}
	}

	ef.updateRouterFactory()
}

// for test.
func (ef *EntryPointFactory) ServerEntryPointsTCP() TCPEntryPoints {
	ef.mu.Lock()
	defer ef.mu.Unlock()

	eps := make(TCPEntryPoints, len(ef.staticEntryPointsTCP)+len(ef.dynamicEntryPointsTCP))
	for key, ep := range ef.staticEntryPointsTCP {
		eps[key] = ep
	}
	for key, ep := range ef.dynamicEntryPointsTCP {
		eps[key] = ep
	}
	return eps
}

// fot test.
func (ef *EntryPointFactory) ServerEntryPointsUDP() UDPEntryPoints {
	ef.mu.Lock()
	defer ef.mu.Unlock()

	eps := make(UDPEntryPoints, len(ef.staticEntryPointsUDP)+len(ef.dynamicEntryPointsUDP))
	for key, ep := range ef.staticEntryPointsUDP {
		eps[key] = ep
	}
	for key, ep := range ef.dynamicEntryPointsUDP {
		eps[key] = ep
	}
	return eps
}

func (ef *EntryPointFactory) EntryPoints() static.EntryPoints {
	ef.mu.Lock()
	defer ef.mu.Unlock()

	results := make(static.EntryPoints, len(ef.staticConfiguration.EntryPoints)+len(ef.dynamicEntryPoints))
	for k, v := range ef.staticConfiguration.EntryPoints {
		results[k] = v
	}
	for k, v := range ef.dynamicEntryPoints {
		results[k] = v
	}
	return results
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

func assignEntryPointTransport(r *dynamic.Router, e *TCPEntryPoint) {
	t := parseEntryPoint(r)
	if t == nil {
		e.transportConfiguration.SetDefaults()
		return
	}
	if httpServer, ok := e.httpServer.Server.(*http.Server); ok {
		httpServer.ReadTimeout = time.Duration(t.RespondingTimeouts.ReadTimeout)
		httpServer.WriteTimeout = time.Duration(t.RespondingTimeouts.WriteTimeout)
		httpServer.IdleTimeout = time.Duration(t.RespondingTimeouts.IdleTimeout)
	}
	if httpServer, ok := e.httpsServer.Server.(*http.Server); ok {
		httpServer.ReadTimeout = time.Duration(t.RespondingTimeouts.ReadTimeout)
		httpServer.WriteTimeout = time.Duration(t.RespondingTimeouts.WriteTimeout)
		httpServer.IdleTimeout = time.Duration(t.RespondingTimeouts.IdleTimeout)
	}
	e.transportConfiguration = t
}

func assignEntryPointTransportForStatic(r *dynamic.Router, e *static.EntryPoint) {
	t := parseEntryPoint(r)
	if t == nil {
		e.Transport.SetDefaults()
		return
	}
	e.Transport = t
}

func parseEntryPoint(r *dynamic.Router) *static.EntryPointsTransport {
	if r.EntryPointTransport == "" {
		return nil
	}

	logger := log.WithoutContext()
	t := &static.EntryPointsTransport{}
	if err := json.Unmarshal([]byte(r.EntryPointTransport), t); err != nil {
		logger.Warn("parse entrypoint transport failed: ", "transport is: ", r.EntryPointTransport, "error is: ", err.Error())
		return nil
	}

	res := &static.EntryPointsTransport{}
	res.SetDefaults()

	if t.RespondingTimeouts.IdleTimeout != 0 {
		res.RespondingTimeouts.IdleTimeout = t.RespondingTimeouts.IdleTimeout
	}

	if t.RespondingTimeouts.ReadTimeout != 0 {
		res.RespondingTimeouts.ReadTimeout = t.RespondingTimeouts.ReadTimeout
	}

	if t.RespondingTimeouts.WriteTimeout != 0 {
		res.RespondingTimeouts.WriteTimeout = t.RespondingTimeouts.WriteTimeout
	}
	return res
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
