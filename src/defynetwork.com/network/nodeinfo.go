package network

import (
	"bytes"
	"io"
	"defynetwork.com/structs"
	"defynetwork.com/tools"
)

func (p *NodeInfo) Pack() []byte {
	w := new(bytes.Buffer)
	p.Dump(w)
	return w.Bytes()
}

func (p *NodeInfo) Dump(w io.Writer) {
	tools.Dump(w, p.Hid)
	tools.Dump(w, uint32(p.Mport))
	tools.Dump(w, uint32(p.Tport))
	tools.Dump(w, uint32(p.Bport))
	tools.Dump(w, tools.Packss(p.Hosts))
}

func (p *NodeInfo) Load(r io.Reader) {
	p.Hid = tools.Loadu64(r)
	p.Mport = int(tools.Loadu32(r))
	p.Tport = int(tools.Loadu32(r))
	p.Bport = int(tools.Loadu32(r))
	p.Hosts = tools.Unpackss(r)
}

type NodeInfo struct {
	Hid uint64
	Hosts []string
	Mport int
	Tport int
	Bport int
}

var TypeNodeInfo = structs.Type.Reg(
	41,
	func(v interface{}) bool {
		_, ok := v.(*NodeInfo)
		return ok
	},
	func(w io.Writer, v interface{}, args ...interface{}) {
		v.(*NodeInfo).Dump(w)
	},
	func(r io.Reader, args ...interface{}) interface{} {
		v := &NodeInfo{}
		v.Load(r)
		return v
	})
