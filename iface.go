package iface

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type IFace struct {
	TLD  string
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (p IFace) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}
	qname := state.Name()

	if !strings.HasSuffix(qname, p.TLD) || (state.QType() != dns.TypeA && state.QType() != dns.TypeAAAA || state.QType() != dns.TypeTXT) {
		return plugin.NextOrFailure(p.Name(), p.Next, ctx, w, r)
	}

	ifaceName := strings.TrimSuffix(qname, p.TLD)

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
		ipPart, maskPart, _ := strings.Cut(addr.String(), "/")
		log.Infof("Replying with address: %s %s", ipPart, maskPart)
		ip := net.ParseIP(ipPart)
		if ip.To4() != nil && state.QType() == dns.TypeA {
			rr := new(dns.A)
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: dns.ClassINET}
			rr.A = ip.To4()

			answers = append(answers, rr)
		} else if ip.To16() != nil && state.QType() == dns.TypeAAAA {
			rr := new(dns.AAAA)
			rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeAAAA, Class: dns.ClassINET}
			rr.AAAA = ip.To16()

			answers = append(answers, rr)
		} else if state.QType() == dns.TypeTXT {
			tt := new(dns.TXT)
			tt.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeTXT, Class: dns.ClassINET}
			tt.Txt = append(tt.Txt, fmt.Sprintf("%s://%s", addr.Network(), addr.String()))
			answers = append(answers, tt)
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
