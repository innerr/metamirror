package network

import (
	"io"
	"testing"
	"defynetwork.com/tools"
)

func TestChannels(t *testing.T) {
	bi := NewBiChannel()
	ac := NewChannels(bi.A)
	a1 := NewChannels(ac.Sub("1"))
	a2 := NewChannels(ac.Sub("2"))
	a3 := NewChannels(a2.Sub("1"))
	bc := NewChannels(bi.B)
	b1 := NewChannels(bc.Sub("1"))
	b2 := NewChannels(bc.Sub("2"))
	b3 := NewChannels(b2.Sub("1"))

	x := uint32(0)
	bc.Receive("30", func(r io.Reader, n uint32) {
		x = 100
	})
	b1.Receive("20", func(r io.Reader, n uint32) {
		x = 101
	})
	b2.Receive("10", func(r io.Reader, n uint32) {
		x = 102
	})
	b3.Receive("00", func(r io.Reader, n uint32) {
		x = 103
	})

	ac.Send("30", nil, 0)
	tools.Check(x, uint32(100))
	a1.Send("20", nil, 0)
	tools.Check(x, uint32(101))
	a2.Send("10", nil, 0)
	tools.Check(x, uint32(102))
	a3.Send("00", nil, 0)
	tools.Check(x, uint32(103))
}

