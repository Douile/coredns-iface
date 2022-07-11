package iface

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.RegisterPlugin("iface", setup) }

func setup(c *caddy.Controller) error {
	c.Next() // 'demo'
	if c.NextArg() {
		return plugin.Error("iface", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return IFace{}
	})

	return nil
}
