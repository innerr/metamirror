package structs

import (
	"io"
)

func (p *AsynChannel) Close() {
	close(p.tasks)
	<-p.done
	p.origin.Close()
}

func (p *AsynChannel) Send(r io.Reader, n uint32) {
	p.tasks <-AsynChannelTask{r, n}
}

func (p *AsynChannel) Receive(fun Transport) {
	p.origin.Receive(fun)
}

func (p *AsynChannel) run() {
	for task := range p.tasks {
		p.origin.Send(task.r, task.n)
	}
	p.done <-true
}

func NewAsynChannel(origin IChannel, backlog int) *AsynChannel {
	if backlog <= 0 {
		backlog = 1024
	}
	p := &AsynChannel{origin, make(chan AsynChannelTask, backlog), make(chan bool)}
	go p.run()
	return p
}

type AsynChannel struct {
	origin IChannel
	tasks chan AsynChannelTask
	done chan bool
}

type AsynChannelTask struct {
	r io.Reader
	n uint32
}

type IChannel interface {
	Send(io.Reader, uint32)
	Receive(Transport)
	Close()
}

type Transport func(io.Reader, uint32)
