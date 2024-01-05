package iface

import (
	"fmt"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin("iface")

func init() { plugin.Register("iface", setup) }

func setup(c *caddy.Controller) error {
	c.Next()

	iface := new(IFace)

	if c.NextArg() {
		iface.TLD = c.Val()

		if !strings.HasSuffix(iface.TLD, ".") {
			iface.TLD = fmt.Sprintf("%s.", iface.TLD)
		}

		if !strings.HasPrefix(iface.TLD, ".") {
			iface.TLD = fmt.Sprintf(".%s", iface.TLD)
		}

	} else {
		iface.TLD = ".iface."
	}

	if c.NextArg() {
		return plugin.Error("iface", c.ArgErr())
	}

	// Add the Plugin to CoreDNS, so Servers can use it in their plugin chain.
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		iface.Next = next

		return iface
	})

	return nil
}
