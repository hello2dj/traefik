package server

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/metrics"
	"github.com/traefik/traefik/v2/pkg/server/middleware"
	"github.com/traefik/traefik/v2/pkg/server/service"
	"github.com/traefik/traefik/v2/pkg/tls"
)

type entpointTestData struct {
	name   string
	conf   dynamic.Configuration
	tcpNum int
	udpNum int
}

func getEntrypointTestData(testServer *httptest.Server) []entpointTestData {
	data := []entpointTestData{
		{
			name:   "initial data 4 tcp, 2 udp",
			tcpNum: 4,
			udpNum: 2,
		}, {
			name:   "remove 1 tcp => 3 tcp, 2 udp",
			tcpNum: 3,
			udpNum: 2,
		}, {
			name:   "remove 1 udp => 3 tcp, 1 udp",
			tcpNum: 3,
			udpNum: 1,
		}, {
			name:   "add 1 tcp => 4 tcp, 1 udp",
			tcpNum: 4,
			udpNum: 1,
		}, {
			name:   "add 1 tcp => 5 tcp, 1 udp",
			tcpNum: 5,
			udpNum: 1,
		}, {
			name:   "remove 1 tcp => 4 tcp, 1 udp",
			tcpNum: 4,
			udpNum: 1,
		}, {
			name:   "add 1 udp => 4 tcp, 2 udp",
			tcpNum: 4,
			udpNum: 2,
		}, {
			name:   "add 1 duplicated tcp entrypoint => 4 tcp, 2 udp",
			tcpNum: 4,
			udpNum: 2,
		}, {
			name:   "add 1 duplicated udp entrypoint => 4 tcp, 2 udp",
			tcpNum: 4,
			udpNum: 2,
		}, {
			name:   "add 1 http router with duplicated entrypoint and a new http router with new entrypoint => 5 tcp, 2 udp",
			tcpNum: 5,
			udpNum: 2,
		},
	}

	for k, v := range entpointFixtures(testServer.URL) {
		data[k].conf = v
	}
	return data
}

func Test(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	staticConfig := static.Configuration{
		EntryPoints: map[string]*static.EntryPoint{
			"web": {},
		},
	}

	roundTripperManager := service.NewRoundTripperManager()
	roundTripperManager.Update(map[string]*dynamic.ServersTransport{"default@internal": {}})
	managerFactory := service.NewManagerFactory(staticConfig, nil, metrics.NewVoidRegistry(), roundTripperManager, nil)
	tlsManager := tls.NewManager()

	for _, e := range staticConfig.EntryPoints {
		e.SetDefaults()
	}

	serverEntryPointsTCP, err := NewTCPEntryPoints(staticConfig.EntryPoints, staticConfig.HostResolver)
	require.NoError(t, err)

	serverEntryPointsUDP, err := NewUDPEntryPoints(staticConfig.EntryPoints)
	require.NoError(t, err)

	factory := NewRouterFactory(staticConfig, managerFactory, tlsManager, middleware.NewChainBuilder(nil, nil, nil), nil, metrics.NewVoidRegistry())

	epFactory := NewEntryPointFactory(factory, staticConfig, serverEntryPointsTCP, serverEntryPointsUDP)

	for _, v := range getEntrypointTestData(testServer) {
		epFactory.BuildEntryPoints(v.conf)
		assert.Equal(t, v.tcpNum, len(epFactory.ServerEntryPointsTCP()), "tcp entrypoints length init")
		assert.Equal(t, v.udpNum, len(epFactory.ServerEntryPointsUDP()), "udp entrypoints length init")
		compareEntryPoints(t, v.conf, epFactory)
	}
}

func compareEntryPoints(t *testing.T, conf dynamic.Configuration, epFactory *EntryPointFactory) {
	var (
		tcpEP       = epFactory.ServerEntryPointsTCP()
		udpEP       = epFactory.ServerEntryPointsUDP()
		entryPoints = map[string]bool{}
	)
	for _, r := range conf.HTTP.Routers {
		for _, e := range r.EntryPoints {
			entryPoints[e] = true
		}
	}

	for _, r := range conf.TCP.Routers {
		for _, e := range r.EntryPoints {
			entryPoints[e] = true
		}
	}

	for _, r := range conf.UDP.Routers {
		for _, e := range r.EntryPoints {
			entryPoints[e] = true
		}
	}

	expectEntryPoints := []string{}
	for e := range entryPoints {
		expectEntryPoints = append(expectEntryPoints, e)
	}

	realEntryPoints := []string{}
	for e := range tcpEP {
		realEntryPoints = append(realEntryPoints, e)
	}

	for e := range udpEP {
		realEntryPoints = append(realEntryPoints, e)
	}

	assert.Equal(t, len(expectEntryPoints), len(realEntryPoints), "entrpoints number is wrong")

	sort.Sort(sort.StringSlice(realEntryPoints))
	sort.Sort(sort.StringSlice(expectEntryPoints))

	for i, v := range expectEntryPoints {
		assert.Equal(t, v, realEntryPoints[i], "entrpoints is wrong")
	}
}
