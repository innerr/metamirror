package tools

import (
	"strconv"
	"strings"
)

func ParseUrl(url string, def int) (host string, port int, err error) {
	host = url
	port = def
	i := strings.LastIndex(url, ":")
	if i >= 0 {
		host = url[:i]
		port, err = strconv.Atoi(url[i + 1:])
	}
	return
}
