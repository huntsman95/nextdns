package main

import (
	"net"
	"strings"
	"testing"

	"github.com/nextdns/nextdns/config"
	"github.com/nextdns/nextdns/discovery"
	"github.com/nextdns/nextdns/resolver/query"
)

func Test_isLocalhostMode(t *testing.T) {
	tests := []struct {
		listens []string
		want    bool
	}{
		{[]string{"127.0.0.1:53"}, true},
		{[]string{"127.0.0.1:5353"}, true},
		{[]string{"10.0.0.1:53"}, false},
		{[]string{"127.0.0.1:53", "10.0.0.1:53"}, false},
		{[]string{"10.0.0.1:53", "127.0.0.1:53"}, false},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.listens, ","), func(t *testing.T) {
			if got := isLocalhostMode(&config.Config{Listens: tt.listens}); got != tt.want {
				t.Errorf("isLocalhostMode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfileNamesForQuery(t *testing.T) {
	mac, err := net.ParseMAC("00:11:22:33:44:55")
	if err != nil {
		t.Fatal(err)
	}
	r := discovery.Resolver{stubDiscoverySource{
		addrs: map[string][]string{
			"10.0.0.2": {"home-pc.local."},
		},
		macs: map[string][]string{
			"00:11:22:33:44:55": {"home-pc.example.lan."},
		},
	}}
	names := profileNamesForQuery(r, query.Query{PeerIP: net.ParseIP("10.0.0.2"), MAC: mac})
	if got, want := len(names), 2; got != want {
		t.Fatalf("len(profileNamesForQuery()) = %d, want %d", got, want)
	}
	if names[0] != "home-pc.local." || names[1] != "home-pc.example.lan." {
		t.Fatalf("profileNamesForQuery() = %v, want both address and MAC names", names)
	}
	if got := profileNamesForQuery(nil, query.Query{PeerIP: net.ParseIP("10.0.0.2"), MAC: mac}); got != nil {
		t.Fatalf("profileNamesForQuery() with nil resolver = %v, want nil", got)
	}
}

type stubDiscoverySource struct {
	addrs map[string][]string
	macs  map[string][]string
}

func (s stubDiscoverySource) Name() string { return "stub" }

func (s stubDiscoverySource) Visit(func(name string, addrs []string)) {}

func (s stubDiscoverySource) LookupAddr(addr string) []string {
	return s.addrs[strings.ToLower(addr)]
}

func (s stubDiscoverySource) LookupHost(name string) []string { return nil }

func (s stubDiscoverySource) LookupMAC(mac string) []string {
	return s.macs[strings.ToLower(mac)]
}
