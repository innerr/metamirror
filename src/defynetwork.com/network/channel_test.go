package network

import (
	"testing"
	"net"
	"defynetwork.com/tools"
)

func TestChannel(t *testing.T) {
	conn, err := net.Dial("tcp", "www.163.com:80")
	if err != nil {
		t.Fatal("conn failed")
	}
	c := NewTcpChannel(tools.NewLog("", true, 0), conn, false)
	go c.Start()
	conn.Close()
}

