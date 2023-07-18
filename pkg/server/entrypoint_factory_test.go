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
	th "github.com/traefik/traefik/v2/pkg/testhelpers"
	"github.com/traefik/traefik/v2/pkg/tls"
)

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

	dynamicConfigs1 := dynamic.Configuration{
		HTTP: th.BuildConfiguration(
			th.WithRouters(
				th.WithRouter("foo",
					th.WithEntryPoints("web"),
					th.WithServiceName("bar"),
					th.WithRule("Path(`/ok`)")),
				th.WithRouter("foo1",
					th.WithEntryPoints("http-30001"),
					th.WithRule("Path(`/unauthorized`)"),
					th.WithServiceName("bar")),
			),
			th.WithLoadBalancerServices(th.WithService("bar",
				th.WithServers(th.WithServer(testServer.URL))),
			),
		),
		TCP: &dynamic.TCPConfiguration{
			Routers: map[string]*dynamic.TCPRouter{
				"Router0": {
					EntryPoints: []string{
						"tcp-30002",
					},
					Service:  "foobar1",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
				"Router1": {
					EntryPoints: []string{
						"tcp-30003",
					},
					Service:  "foobar2",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
			},
			Services: map[string]*dynamic.TCPService{
				"foobar1": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
				"foobar2": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
			},
		},
		UDP: &dynamic.UDPConfiguration{
			Routers: map[string]*dynamic.UDPRouter{
				"Router0": {
					EntryPoints: []string{
						"udp-30004",
					},
					Service: "fiibar1",
				},
				"Router1": {
					EntryPoints: []string{
						"udp-30005",
					},
					Service: "fiibar2",
				},
			},
			Services: map[string]*dynamic.UDPService{
				"fiibar1": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
				"fiibar2": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
			},
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
	epFactory.BuildEntryPoints(dynamicConfigs1)
	assert.Equal(t, 4, len(epFactory.ServerEntryPointsTCP()), "tcp entrypoints length init")
	assert.Equal(t, 2, len(epFactory.ServerEntryPointsUDP()), "udp entrypoints length init")
	compareEntryPoints(t, dynamicConfigs1, epFactory)

	dynamicConfigs2 := dynamic.Configuration{
		HTTP: th.BuildConfiguration(
			th.WithRouters(
				th.WithRouter("foo",
					th.WithEntryPoints("web"),
					th.WithRule("Path(`/unauthorized`)"),
					th.WithServiceName("bar")),
			),
			th.WithLoadBalancerServices(th.WithService("bar",
				th.WithServers(th.WithServer(testServer.URL))),
			),
		),
		TCP: &dynamic.TCPConfiguration{
			Routers: map[string]*dynamic.TCPRouter{
				"Router0": {
					EntryPoints: []string{
						"tcp-30002",
					},
					Service:  "foobar1",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
				"Router1": {
					EntryPoints: []string{
						"tcp-30003",
					},
					Service:  "foobar2",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
			},
			Services: map[string]*dynamic.TCPService{
				"foobar1": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
				"foobar2": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
			},
		},
		UDP: &dynamic.UDPConfiguration{
			Routers: map[string]*dynamic.UDPRouter{
				"Router0": {
					EntryPoints: []string{
						"udp-30004",
					},
					Service: "fiibar1",
				},
				"Router1": {
					EntryPoints: []string{
						"udp-30005",
					},
					Service: "fiibar2",
				},
			},
			Services: map[string]*dynamic.UDPService{
				"fiibar1": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
				"fiibar2": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
			},
		},
	}
	epFactory.BuildEntryPoints(dynamicConfigs2)
	assert.Equal(t, 3, len(epFactory.ServerEntryPointsTCP()), "tcp entrypoints length remove http router")
	assert.Equal(t, 2, len(epFactory.ServerEntryPointsUDP()), "udp entrypoints length remove http router")
	compareEntryPoints(t, dynamicConfigs2, epFactory)

	dynamicConfigs3 := dynamic.Configuration{
		HTTP: th.BuildConfiguration(
			th.WithRouters(
				th.WithRouter("foo",
					th.WithEntryPoints("web"),
					th.WithRule("Path(`/unauthorized`)"),
					th.WithServiceName("bar")),
			),
			th.WithLoadBalancerServices(th.WithService("bar",
				th.WithServers(th.WithServer(testServer.URL))),
			),
		),
		TCP: &dynamic.TCPConfiguration{
			Routers: map[string]*dynamic.TCPRouter{
				"Router0": {
					EntryPoints: []string{
						"tcp-30002",
					},
					Service:  "foobar1",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
				"Router1": {
					EntryPoints: []string{
						"tcp-30003",
					},
					Service:  "foobar2",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
			},
			Services: map[string]*dynamic.TCPService{
				"foobar1": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
				"foobar2": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
			},
		},
		UDP: &dynamic.UDPConfiguration{
			Routers: map[string]*dynamic.UDPRouter{
				"Router0": {
					EntryPoints: []string{
						"udp-30004",
					},
					Service: "fiibar1",
				},
			},
			Services: map[string]*dynamic.UDPService{
				"fiibar1": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
				"fiibar2": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
			},
		},
	}
	epFactory.BuildEntryPoints(dynamicConfigs3)
	assert.Equal(t, 3, len(epFactory.ServerEntryPointsTCP()), "tcp entrypoints length remove udp")
	assert.Equal(t, 1, len(epFactory.ServerEntryPointsUDP()), "udp entrypoints length remove udp")
	compareEntryPoints(t, dynamicConfigs3, epFactory)

	dynamicConfigs4 := dynamic.Configuration{
		HTTP: th.BuildConfiguration(
			th.WithRouters(
				th.WithRouter("foo",
					th.WithEntryPoints("web"),
					th.WithRule("Path(`/unauthorized`)"),
					th.WithServiceName("bar")),
			),
			th.WithLoadBalancerServices(th.WithService("bar",
				th.WithServers(th.WithServer(testServer.URL))),
			),
		),
		TCP: &dynamic.TCPConfiguration{
			Routers: map[string]*dynamic.TCPRouter{
				"Router0": {
					EntryPoints: []string{
						"tcp-30002",
					},
					Service:  "foobar1",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
				"Router1": {
					EntryPoints: []string{
						"tcp-30003",
					},
					Service:  "foobar2",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
				"Router2": {
					EntryPoints: []string{
						"tcp-30004",
					},
					Service:  "foobar2",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
			},
			Services: map[string]*dynamic.TCPService{
				"foobar1": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
				"foobar2": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
			},
		},
		UDP: &dynamic.UDPConfiguration{
			Routers: map[string]*dynamic.UDPRouter{
				"Router0": {
					EntryPoints: []string{
						"udp-30004",
					},
					Service: "fiibar1",
				},
			},
			Services: map[string]*dynamic.UDPService{
				"fiibar1": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
				"fiibar2": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
			},
		},
	}

	epFactory.BuildEntryPoints(dynamicConfigs4)
	assert.Equal(t, 4, len(epFactory.ServerEntryPointsTCP()), "tcp entrypoints length add tcp")
	assert.Equal(t, 1, len(epFactory.ServerEntryPointsUDP()), "udp entrypoints length add tcp")
	compareEntryPoints(t, dynamicConfigs4, epFactory)

	dynamicConfigs5 := dynamic.Configuration{
		HTTP: th.BuildConfiguration(
			th.WithRouters(
				th.WithRouter("foo",
					th.WithEntryPoints("web"),
					th.WithRule("Path(`/unauthorized`)"),
					th.WithServiceName("bar")),
				th.WithRouter("foo1",
					th.WithEntryPoints("http-31001"),
					th.WithRule("Path(`/unauthorized`)"),
					th.WithServiceName("bar")),
			),
			th.WithLoadBalancerServices(th.WithService("bar",
				th.WithServers(th.WithServer(testServer.URL))),
			),
		),
		TCP: &dynamic.TCPConfiguration{
			Routers: map[string]*dynamic.TCPRouter{
				"Router0": {
					EntryPoints: []string{
						"tcp-30002",
					},
					Service:  "foobar1",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
				"Router1": {
					EntryPoints: []string{
						"tcp-30003",
					},
					Service:  "foobar2",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
				"Router2": {
					EntryPoints: []string{
						"tcp-30004",
					},
					Service:  "foobar2",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
			},
			Services: map[string]*dynamic.TCPService{
				"foobar1": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
				"foobar2": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
			},
		},
		UDP: &dynamic.UDPConfiguration{
			Routers: map[string]*dynamic.UDPRouter{
				"Router0": {
					EntryPoints: []string{
						"udp-30004",
					},
					Service: "fiibar1",
				},
			},
			Services: map[string]*dynamic.UDPService{
				"fiibar1": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
				"fiibar2": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
			},
		},
	}
	epFactory.BuildEntryPoints(dynamicConfigs5)
	assert.Equal(t, 5, len(epFactory.ServerEntryPointsTCP()), "tcp entrypoints length add http router")
	assert.Equal(t, 1, len(epFactory.ServerEntryPointsUDP()), "udp entrypoints length add http router")
	compareEntryPoints(t, dynamicConfigs5, epFactory)

	dynamicConfigs6 := dynamic.Configuration{
		HTTP: th.BuildConfiguration(
			th.WithRouters(
				th.WithRouter("foo",
					th.WithEntryPoints("web"),
					th.WithRule("Path(`/unauthorized`)"),
					th.WithServiceName("bar")),
				th.WithRouter("foo1",
					th.WithEntryPoints("http-31001"),
					th.WithRule("Path(`/unauthorized`)"),
					th.WithServiceName("bar")),
			),
			th.WithLoadBalancerServices(th.WithService("bar",
				th.WithServers(th.WithServer(testServer.URL))),
			),
		),
		TCP: &dynamic.TCPConfiguration{
			Routers: map[string]*dynamic.TCPRouter{
				"Router0": {
					EntryPoints: []string{
						"tcp-30002",
					},
					Service:  "foobar1",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
				"Router2": {
					EntryPoints: []string{
						"tcp-30004",
					},
					Service:  "foobar2",
					Rule:     "foobar",
					Priority: 42,
					TLS: &dynamic.RouterTCPTLSConfig{
						Passthrough: false,
						Options:     "foo",
					},
				},
			},
			Services: map[string]*dynamic.TCPService{
				"foobar1": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
				"foobar2": {
					LoadBalancer: &dynamic.TCPServersLoadBalancer{
						Servers: []dynamic.TCPServer{
							{
								Port: "42",
							},
						},
						TerminationDelay: func(i int) *int { return &i }(42),
						ProxyProtocol:    &dynamic.ProxyProtocol{Version: 42},
					},
				},
			},
		},
		UDP: &dynamic.UDPConfiguration{
			Routers: map[string]*dynamic.UDPRouter{
				"Router0": {
					EntryPoints: []string{
						"udp-30004",
					},
					Service: "fiibar1",
				},
			},
			Services: map[string]*dynamic.UDPService{
				"fiibar1": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
				"fiibar2": {
					LoadBalancer: &dynamic.UDPServersLoadBalancer{
						Servers: []dynamic.UDPServer{
							{
								Port: "42",
							},
						},
					},
				},
			},
		},
	}

	epFactory.BuildEntryPoints(dynamicConfigs6)
	assert.Equal(t, 4, len(epFactory.ServerEntryPointsTCP()), "tcp entrypoints length remove tcp again")
	assert.Equal(t, 1, len(epFactory.ServerEntryPointsUDP()), "udp entrypoints length remove tcp again")
	compareEntryPoints(t, dynamicConfigs6, epFactory)
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
