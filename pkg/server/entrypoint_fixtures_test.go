package server

import (
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	th "github.com/traefik/traefik/v2/pkg/testhelpers"
)

func entpointFixtures(url string) []dynamic.Configuration {
	return []dynamic.Configuration{
		{
			HTTP: th.BuildConfiguration(
				th.WithRouters(
					th.WithRouter("foo",
						th.WithEntryPoints("web"),
						th.WithServiceName("bar"),
						th.WithRule("Path(`/ok`)"),
						th.WithEntryPointTransport(`{"respondingTimeouts":{"readTimeout":3,"writeTimeout":4}}`),
					),
					th.WithRouter("foo1",
						th.WithEntryPoints("http-30001"),
						th.WithRule("Path(`/unauthorized`)"),
						th.WithServiceName("bar"),
						th.WithEntryPointTransport(`{"respondingTimeouts":{"readTimeout":3,"writeTimeout":4}}`),
					),
				),
				th.WithLoadBalancerServices(th.WithService("bar",
					th.WithServers(th.WithServer(url))),
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
		},
		{
			HTTP: th.BuildConfiguration(
				th.WithRouters(
					th.WithRouter("foo",
						th.WithEntryPoints("web"),
						th.WithRule("Path(`/unauthorized`)"),
						th.WithServiceName("bar")),
				),
				th.WithLoadBalancerServices(th.WithService("bar",
					th.WithServers(th.WithServer(url))),
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
		},
		{
			HTTP: th.BuildConfiguration(
				th.WithRouters(
					th.WithRouter("foo",
						th.WithEntryPoints("web"),
						th.WithRule("Path(`/unauthorized`)"),
						th.WithServiceName("bar")),
				),
				th.WithLoadBalancerServices(th.WithService("bar",
					th.WithServers(th.WithServer(url))),
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
		},
		{
			HTTP: th.BuildConfiguration(
				th.WithRouters(
					th.WithRouter("foo",
						th.WithEntryPoints("web"),
						th.WithRule("Path(`/unauthorized`)"),
						th.WithServiceName("bar")),
				),
				th.WithLoadBalancerServices(th.WithService("bar",
					th.WithServers(th.WithServer(url))),
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
		},
		{
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
					th.WithServers(th.WithServer(url))),
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
		},
		{
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
					th.WithServers(th.WithServer(url))),
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
		},
		{
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
					th.WithServers(th.WithServer(url))),
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
					"Router1": {
						EntryPoints: []string{
							"udp-30007",
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
		},
		{
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
					th.WithServers(th.WithServer(url))),
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
					"Router3": {
						EntryPoints: []string{
							"tcp-30004",
						},
						Service:  "foobar2",
						Rule:     "foobarxxxxxxx",
						Priority: 42,
						TLS: &dynamic.RouterTCPTLSConfig{
							Passthrough: false,
							Options:     "fooxxxxxxx",
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
							"udp-30007",
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
		},
		{
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
					th.WithServers(th.WithServer(url))),
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
					"Router3": {
						EntryPoints: []string{
							"tcp-30004",
						},
						Service:  "foobar2",
						Rule:     "foobarxxxxxxx",
						Priority: 42,
						TLS: &dynamic.RouterTCPTLSConfig{
							Passthrough: false,
							Options:     "fooxxxxxxx",
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
							"udp-30007",
						},
						Service: "fiibar2",
					},
					"Router2": {
						EntryPoints: []string{
							"udp-30007",
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
		},
		{
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
					th.WithRouter("foo2",
						th.WithEntryPoints("http-31001"),
						th.WithRule("Path(`/unauthorizedxxx`)"),
						th.WithServiceName("bar")),
					th.WithRouter("foo3",
						th.WithEntryPoints("http-31002"),
						th.WithRule("Path(`/unauthorizedyyyyy`)"),
						th.WithServiceName("bar")),
				),
				th.WithLoadBalancerServices(th.WithService("bar",
					th.WithServers(th.WithServer(url))),
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
					"Router3": {
						EntryPoints: []string{
							"tcp-30004",
						},
						Service:  "foobar2",
						Rule:     "foobarxxxxxxx",
						Priority: 42,
						TLS: &dynamic.RouterTCPTLSConfig{
							Passthrough: false,
							Options:     "fooxxxxxxx",
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
							"udp-30007",
						},
						Service: "fiibar2",
					},
					"Router2": {
						EntryPoints: []string{
							"udp-30007",
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
		},
	}
}
