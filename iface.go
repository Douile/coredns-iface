package iface

import (
	"context"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type IFace struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (p IFace) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	if !strings.HasSuffix(qname, ".iface.") || (state.QType() != dns.TypeA && state.QType() != dns.TypeAAAA) {
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}

	ifaceName := strings.TrimSuffix(qname, ".iface.")

	iface, err := net.InterfaceByName(ifaceName)

	if err != nil {
		return dns.RcodeNameError, nil
	}

	addrs, err := iface.Addrs()

	if err != nil {
		return dns.RcodeNameError, nil
	}

	answers := []dns.RR{}
	for _, addr := range addrs {
		ip := net.ParseIP(addr.String())
		if ip.To4() != nil {
			rr := new(dns.A)
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: dns.ClassINET}
			rr.A = ip.To4()

			answers = append(answers, rr)
		} else if ip.To16() != nil {
			rr := new(dns.AAAA)
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeAAAA, Class: dns.ClassINET}
			rr.AAAA = ip.To16()

			answers = append(answers, rr)
		}
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (p IFace) Name() string { return "iface" }
