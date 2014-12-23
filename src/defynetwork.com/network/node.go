package network

import (
	"net"
	"strings"
	"defynetwork.com/structs"
	"defynetwork.com/tools"
)

func (p *Node) SearchLan(interval int) {
	if p.Info.Bport <= 0 {
		return
	}

	log := p.log.Mod("broadcast")
	go func() {
		defer tools.Catch(log.Detail, false)
		p.Broadcastor.Ping(interval)
	}()

	p.Broadcastor.Serve(func(host string, info *NodeInfo) {
		go func() {
			defer tools.Catch(log.Debug, true)
			p.Conn(host, info.Mport)
		}()
	})
}

func (p *Node) bind(conn net.Conn, chs *Channels, passive bool) {
	host := strings.Split(conn.RemoteAddr().String(), ":")[0]

	p.Conns.Bind(conn, chs, host, passive, func() {
		if !passive && p.Transfer.ReqCount() != 0 {
			p.Transfer.Conn(p.Conns.ByMain(chs))
		}
	})

	p.Domains.Bind(chs, passive)
}

func (p *Node) unbind(conn net.Conn, chs *Channels) {
	info := p.Conns.ByMain(chs)
	if info != nil {
		p.Neighbors.Offline(info.Hid)
	}
	p.Domains.Unbind(chs)
	p.Conns.Unbind(conn, chs)
}

func NewNode(log *tools.Log, hid uint64, mport, tport, bport int, hashs IHashLib, newdm NewDmFunc) *Node {
	addrs := NewMyAddrs(log.Mod("addrs"))
	info := &NodeInfo{hid, addrs.MyAddrs(false), mport, tport, bport}
	conns := NewConnections(log.Mod("conns"), info)

	// BUG: lazy mode
	domains := NewDomains(log.Mod("domains"), newdm, false)

	neighbors := NewNeighbors(log.Mod("neighbors"), info)
	flags := &DomainFlags{true, true, false, true, true, true}
	domain := NewDomain(log.Mod("dm:neighbors!"), neighbors, &structs.DumbPersist{}, flags)
	domains.Reg("neighbors!", domain)
	transfer := NewTransfer(log.Mod("transfer"), addrs, tport, hashs, conns, 10000)
	bc := NewBroadcastor(log.Mod("broadcast"), addrs, info)
	p := &Node{log, nil, bc, domains, transfer, conns, neighbors, info, addrs}
	p.TcpNode = NewTcpNode(log, addrs, mport, false, p.bind, p.unbind)
	return p
}

type Node struct {
	log *tools.Log
	*TcpNode
	Broadcastor *Broadcastor
	Domains *Domains
	Transfer * Transfer
	Conns *Connections
	Neighbors *Neighbors
	Info *NodeInfo
	addrs *MyAddrs
}
