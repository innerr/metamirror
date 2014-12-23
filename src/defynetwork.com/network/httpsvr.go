package network

import (
	"bytes"
	"net/http"
	"strconv"
	"fmt"
	"defynetwork.com/structs"
	"defynetwork.com/tools"
)

func (p *HttpSvr) Serve(port int, debug bool) {
	sync := func(w http.ResponseWriter, req *http.Request) {
		r := req.Body
		c := tools.Loadu16(r)
		buf := new(bytes.Buffer)
		tools.Dump(buf, c)
		for i := uint16(0); i < c; i++ {
			name := tools.Loads(r)
			clocks := structs.Clocks{}
			clocks.Load(r)
			domain := p.domains.Get(name)
			if domain == nil {
				panic("domain not found")
			}
			tools.Dumps(buf, name)
			domain.Clocks().Dump(buf)
			delta := domain.Delta(clocks)
			delta.Dump(buf)
		}
		w.Write(buf.Bytes())
	}

	complete := func(w http.ResponseWriter, req *http.Request) {
		r := req.Body
		c := tools.Loadu16(r)
		for i := uint16(0); i < c; i++ {
			name := tools.Loads(r)
			delta := structs.Delta{}
			delta.Load(r)
			domain := p.domains.Get(name)
			if domain == nil {
				panic("domain not found")
			}
			domain.Merge(delta)
		}
		w.Write([]byte("ok"))
	}

	handle := func(w http.ResponseWriter, req *http.Request, h func(http.ResponseWriter, *http.Request)) {
		defer func() {
			if debug {
				return
			}
			if err := recover(); err != nil {
				println(fmt.Sprintf("%v", err))
			}
		}()
		h(w, req)
	}

	http.HandleFunc("/sync", func(w http.ResponseWriter, req *http.Request) {
		handle(w, req, sync)
	})
	http.HandleFunc("/complete", func(w http.ResponseWriter, req *http.Request) {
		handle(w, req, complete)
	})
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}

func NewHttpSvr(domains *Domains) *HttpSvr {
	return &HttpSvr{domains, nil}
}

type HttpSvr struct {
	domains *Domains
	conns []*TcpChannel
}

func (p *HttpClient) Start(host string, port int) {
	addr := host + ":" + strconv.Itoa(port)
	buf := new(bytes.Buffer)
	names := p.domains.Names()
	tools.Dump(buf, uint16(len(names)))
	for _, name := range names {
		tools.Dumps(buf, name)
		domain := p.domains.Get(name)
		domain.Clocks().Dump(buf)
	}
	resp, err := http.Post(addr + "/sync", "application/octet-stream", bytes.NewReader(buf.Bytes()))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	r := resp.Body
	c := tools.Loadu16(r)
	buf = new(bytes.Buffer)
	tools.Dump(buf, c)
	for i := uint16(0); i < c; i++ {
		name := tools.Loads(r)
		clocks := structs.Clocks{}
		clocks.Load(r)
		delta := structs.Delta{}
		delta.Load(r)
		domain := p.domains.Get(name)
		if domain != nil {
			panic("domain not found")
		}
		domain.Merge(delta)
		complete := domain.Delta(clocks)
		tools.Dumps(buf, name)
		complete.Dump(buf)
	}
	resp, err = http.Post(addr + "/complete", "application/octet-stream", bytes.NewReader(buf.Bytes()))
	if err != nil {
		panic(err)
	}
}

func NewHttpClient(domains Domains) *HttpClient {
	return &HttpClient{domains}
}

type HttpClient struct {
	domains Domains
}

// LATER httpsvr
