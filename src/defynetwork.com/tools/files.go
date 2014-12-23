package tools

import (
	"os"
	"path/filepath"
)

func (p *Files) Close() {
	for _, file := range p.files {
		file.Close()
	}
}

func (p *Files) Open(rel string) *os.File {
	path := p.root + "/" + rel
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0777)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	file, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	p.files = append(p.files, file)
	return file
}

func NewFiles(root string) *Files {
	return &Files{root, nil}
}

type Files struct {
	root string
	files []*os.File
}
