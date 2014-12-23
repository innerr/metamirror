package network

import (
	"fmt"
	"defynetwork.com/structs"
	"defynetwork.com/tools"
)

func (p *Neighbors) Count() int {
	size := 0
	p.Dict.Walk(func(path []string, parent *structs.Dict, val structs.DictVal) {
		if !val.Deleted {
			size += 1
		}
	})
	return size
}

func (p *Neighbors) Walk(fun func(*NodeInfo)) {
	p.Dict.Walk(func(path []string, parent *structs.Dict, val structs.DictVal) {
		if val.Deleted {
			return
		}
		neighbor, ok := val.Data.(*NodeInfo)
		if ok {
			fun(neighbor)
		}
	})
}

func (p *Neighbors) Offline(hid uint64) {
	idstr := fmt.Sprintf("%x", hid)
	p.Dict.Del(idstr, true)
	p.log.Debug("lan nodes: ", p.Dict.Size(false), ", offline: ", idstr)
}

func NewNeighbors(log *tools.Log, info *NodeInfo) *Neighbors {
	idstr := fmt.Sprintf("%x", info.Hid)
	changed := func(d *structs.Dict, key string, val structs.DictVal) {
		if idstr == key {
			d.Del(idstr, false)
		}
		log.Debug("lan nodes: ", d.Size(false), ", " + map[bool]string{true: "off", false: "on"}[val.Deleted] + "line: ", key)
	}
	dict := structs.NewDict(log.Mod("dict:neighbors!"), info.Hid)
	dict.SetEventHandle(changed, nil, nil)
	dict.Set(idstr, info)
	return &Neighbors{dict, log}
}

type Neighbors struct {
	*structs.Dict
	log *tools.Log
}
