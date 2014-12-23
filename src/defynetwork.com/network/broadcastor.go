package network

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"time"
	"net"
	"defynetwork.com/tools"
)

func (p *Broadcastor) StopPing() {
	p.ping = false
}

func (p *Broadcastor) Ping(interval int) {
	p.ping = true
	buf := new(bytes.Buffer)
	p.info.Dump(buf)
	data := buf.Bytes()

	for _ = range time.Tick(time.Duration(interval) * time.Millisecond) {
		if !p.ping {
			break
		}
		if p.conn == nil {
			continue
		}
		for addr, _ := range p.baddrs {
			_, err := p.conn.WriteToUDP(data, addr)
			if err != nil {
				p.log.Detail("error: ", err)
			}
		}
	}
}

func (p *Broadcastor) Serve(receive func(host string, info *NodeInfo)) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: p.info.Bport})
	if err != nil {
		conn, err = net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: p.info.Bport + 1})
		if err != nil {
			panic(err)
		} else {
			p.log.Debug("listen: ", p.info.Bport + 1)
		}
	} else {
		p.log.Debug("listen: ", p.info.Bport)
	}
	p.conn = conn

	data := make([]byte, 2048)
	for {
		read, remote, err := conn.ReadFromUDP(data)
		if err != nil {
			p.log.Detail("error: ", err.Error())
			continue
		}
		host := strings.Split(remote.String(), ":")[0]
		buf := bytes.NewReader(data[:read])
		info := &NodeInfo{}
		info.Load(buf)
		receive(host, info)
	}
}

func (p *Broadcastor) UpdateAddrs() {
	baddrs := make(map[*net.UDPAddr]bool)
	for _, host := range p.addrs.MyAddrs(false) {
		addr1, _ := net.ResolveUDPAddr("udp", host + ":" + strconv.Itoa(p.info.Bport))
		if addr1 != nil {
			baddrs[addr1] = true
		}
		addr2, _ := net.ResolveUDPAddr("udp", host + ":" + strconv.Itoa(p.info.Bport + 1))
		if addr1 != nil {
			baddrs[addr2] = true
		}
	}
	if !reflect.DeepEqual(p.baddrs, baddrs) {
		for addr, _ := range baddrs {
			p.log.Msg("broad addr: ", addr)
		}
	}
	p.baddrs = baddrs
}

func NewBroadcastor(log *tools.Log, addrs *MyAddrs, info *NodeInfo) *Broadcastor {
	p := &Broadcastor{log, addrs, info, nil, nil, false}
	p.UpdateAddrs()
	return p
}

type Broadcastor struct {
	log *tools.Log
	addrs *MyAddrs
	info *NodeInfo
	conn *net.UDPConn
	baddrs map[*net.UDPAddr]bool
	ping bool
}
