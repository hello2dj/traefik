package types

import "github.com/traefik/traefik/v2/pkg/config/static"

// EntrypointsGetter avoid importing cycle
type EntrypointsGetter interface {
	EntryPoints() static.EntryPoints
}
