package iface

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("iface", setup) }

func setup(c *caddy.Controller) error {
	c.Next() // 'demo'
	if c.NextArg() {
		return plugin.Error("iface", c.ArgErr())
	}


  iface := new(IFace)

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		iface.Next = next

		return iface
	})

	return nil
}
