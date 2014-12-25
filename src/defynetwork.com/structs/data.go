package structs

type IData interface {
	Pack() []byte
	Commit() []byte
	Merge([]byte)
}
