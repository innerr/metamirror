package tools

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func (p *Log) Detail(msg ...interface{}) {
	if p.debug <= LogLvlDetail {
		p.write(LogLvlDetail, msg...)
	}
}

func (p *Log) Debug(msg ...interface{}) {
	if p.debug <= LogLvlDebug {
		p.write(LogLvlDebug, msg...)
	}
}

func (p *Log) Msg(msg ...interface{}) {
	if p.debug <= LogLvlNormal {
		p.write(LogLvlNormal, msg...)
	}
}

func (p *Log) Close() {
	if p.file != nil {
		p.file.Close()
	}
}

func (p *Log) write(debug int, msg ...interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	mod := " "
	if p.mod != "" {
		mod += "[" + p.mod + "] "
	}
	line := time.Now().Format("2006/01/02 15:04:05") + " " + fmt.Sprint(debug) + mod + fmt.Sprint(msg...)
	if p.screen {
		println(line)
	}
	if p.file == nil {
		return
	}
	p.file.Write([]byte(line + "\n"))
}

func (p *Log) Mod(mod string) *Log {
	return &Log{p.file, mod, p.screen, p.debug, sync.Mutex{}}
}

func NewLog(path string, screen bool, debug int) *Log {
	file := (*os.File)(nil)
	if path != "" {
		size := int64(0)
		info, err := os.Stat(path)
		if err == nil {
			size = info.Size()
		}
		file, err = os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		if size > 1024 * 1024 * 32 {
			file.Truncate(0)
		}
	}
	return &Log{file, "", screen, debug, sync.Mutex{}}
}

type Log struct {
	file *os.File
	mod string
	screen bool
	debug int
	lock sync.Mutex
}

const (
	LogLvlDetail = 0
	LogLvlDebug = 1
	LogLvlNormal = 2
	LogLvlNone = 3
)
