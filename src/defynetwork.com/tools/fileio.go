package tools

import (
	"bytes"
	"crypto/sha1"
	"io"
	"os"
	"path/filepath"
)

func NotExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true
	}
	return false
}

func Cat(data ...[]byte) []byte {
	buf := new(bytes.Buffer)
	for _, it := range data {
		_, err := buf.Write(it)
		if err != nil {
			panic(err)
		}
	}
	return buf.Bytes()
}

func ReadN(r io.Reader, n uint32) []byte {
	buf := new(bytes.Buffer)
	_, err := io.CopyN(buf, r, int64(n))
	if err != nil {
		panic(err)
	}
	data := buf.Bytes()
	return data
}

func OpenW(path string, new bool) (*os.File, error) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0777)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	flag := os.O_RDWR | os.O_CREATE
	if new {
		flag = flag | os.O_TRUNC
	}
	return os.OpenFile(path, flag, 0666)
}

func Mkdir(path string) error {
	return os.MkdirAll(path, 0777)
}

func FileSha1Ex(path string, log *Log) (hash []byte, unexist bool) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, true
		}
		if log != nil {
			log.Msg("error: ", err)
		}
		return nil, false
	}
	defer file.Close()
	sha1 := sha1.New()
	_, err = io.Copy(sha1, file)
	if err != nil {
		if log != nil {
			log.Msg("error: ", err)
		}
		return nil, false
	}
	return sha1.Sum(nil), false
}

func FileSha1(path string, log *Log) []byte {
	hash, _ := FileSha1Ex(path, log)
	return hash
}
