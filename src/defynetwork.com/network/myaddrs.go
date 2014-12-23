package network

import (
	"reflect"
	"strings"
	"net"
	"defynetwork.com/tools"
)

func (p *MyAddrs) MyAddrs(loopback bool) []string {
	addrs := []string{}
	if loopback {
		addrs = []string{"localhost", "127.0.0.1"}
	}
	for addr, _ := range p.ips {
		addrs = append(addrs, addr)
	}
	return addrs
}

func (p *MyAddrs) Update() {
	ips := make(map[string]bool)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		ip := strings.Split(addr.String(), "/")[0]
		if strings.Index(ip, ":") >= 0 {
			continue
		}
		if strings.Index(ip, ".") < 0 {
			continue
		}
		if tools.Loopback(ip) {
			continue
		}
		if tools.NetworkAddr(ip) {
			continue
		}
		ips[ip] = true
	}

	hosts := []string{}
	if len(ips) == 0 {
		hosts = append(hosts, "255.255.255.255")
	} else {
		for ip, _ := range ips {
			host := ip[:strings.LastIndex(ip, ".") + 1] + "255"
			hosts = append(hosts, host)
		}
	}

	if !reflect.DeepEqual(p.ips, ips) {
		for addr, _ := range ips {
			p.log.Msg("my addr: ", addr)
		}
	}
	p.ips = ips
}

func NewMyAddrs(log *tools.Log) *MyAddrs {
	p := &MyAddrs{log, nil}
	p.Update()
	return p
}

type MyAddrs struct {
	log *tools.Log
	ips map[string]bool
}
