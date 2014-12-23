package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"compress/zlib"
	"crypto/rand"
	"io"
	"io/ioutil"
)

func Decode(data []byte, keyword []byte) []byte {
	flag := data[0]
	data = data[1:]

	if flag & _CodedDataEncrypt != 0 {
		iv := data[:aes.BlockSize]
		data = data[aes.BlockSize:]

		block, err := aes.NewCipher(keyword)
		if err != nil {
			panic(err)
		}
		stream := cipher.NewCFBDecrypter(block, iv[:])

		eb := new(bytes.Buffer)
		w := &cipher.StreamWriter{S: stream, W: NewWriteCloser(eb)}
		_, err = io.CopyN(w, bytes.NewReader(data), int64(len(data)))
		if err != nil {
			panic(err)
		}
		w.Close()
		data = eb.Bytes()
	}

	if flag & _CodedDataCompress != 0 {
		buf, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		data, err = ioutil.ReadAll(buf)
		if err != nil {
			panic(err)
		}
	}

	return data
}

func Encode(r io.Reader, n uint32, compress bool, keyword []byte) []byte {
	flag := byte(0)
	if len(keyword) > 0 {
		flag |= _CodedDataEncrypt
	}
	if compress {
		flag |= _CodedDataCompress
	}

	buf := new(bytes.Buffer)
	w := io.Writer(buf)
	c := []io.Closer{}

	_, err := w.Write([]byte{flag})
	if err != nil {
		panic(err)
	}

	if len(keyword) > 0 {
		var iv [aes.BlockSize]byte
		if _, err := io.ReadFull(rand.Reader, iv[:]); err != nil {
			panic(err)
		}
		block, err := aes.NewCipher(keyword)
		if err != nil {
			panic(err)
		}
		stream := cipher.NewCFBEncrypter(block, iv[:])

		_, err = w.Write(iv[:])
		if err != nil {
			panic(err)
		}

		cw := &cipher.StreamWriter{S: stream, W: NewWriteCloser(w)}
		c = append(c, cw)
		w = cw
	}

	if compress {
		zw := zlib.NewWriter(w)
		c = append(c, zw)
		w = zw
	}

	_, err = io.CopyN(w, r, int64(n))
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(c); i++ {
		c[len(c) - i - 1].Close()
	}
	return buf.Bytes()
}

func (p *WriteCloser) Write(data []byte) (int, error) {
	return p.w.Write(data)
}

func (p *WriteCloser) Close() error {
	return nil
}

func NewWriteCloser(w io.Writer) io.WriteCloser {
	return &WriteCloser{w}
}

type WriteCloser struct {
	w io.Writer
}

const (
	_CodedDataCompress = byte(1)
	_CodedDataEncrypt = byte(2)
)
