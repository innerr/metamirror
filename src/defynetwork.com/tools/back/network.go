package tools

import "strings"

func NetworkAddr(host string) bool {
	for _, it := range _NetworkAddrs {
		if it == host {
			return true
		}
	}
	return false
}

var _NetworkAddrs = []string {
	"",
	"0.0.0.0",
}

func Loopback(host string) bool {
	for _, it := range Loopbacks {
		if it == host {
			return true
		}
	}
	return false
}

var Loopbacks = []string{
	"127.0.0.1",
	"localhost",
}

func NetworkErr(err error) bool {
	if err == nil {
		return false
	}
	m := err.Error()
	if _, ok := _NetworkErrs[m]; ok {
		return true
	}
	for _, suffix := range _NetworkErrSuffixs {
		if strings.HasSuffix(m, suffix) {
			return true
		}
	}
	for _, term := range _NetworkErrTerms {
		if strings.Index(m, term) >= 0 {
			return true
		}
	}
	return false
}

var _NetworkErrs = map[string]bool{
	"use of closed network connection": true,
}

var _NetworkErrSuffixs = []string{
	"closed pipe",
	"broken pipe",
	"connection reset by peer",
	"connection refused",
	"no route to host",
	"no route to host [recovered]",
	"operation timed out",
}

var _NetworkErrTerms = []string{
	"connection was aborted",
	"no such host",
	"use of closed network connection",
	"connection refused",
}
