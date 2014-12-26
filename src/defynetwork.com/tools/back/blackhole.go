package tools

import "io"

func Eat(log *Log, r io.Reader, n uint32) {
	_, err := io.CopyN(&BlackHole{}, r, int64(n))
	if err != nil {
		log.Msg("error: ", err)
	}
}

func (p BlackHole) Write(data []byte) (int, error) {
	return len(data), nil
}

type BlackHole struct {
}

