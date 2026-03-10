package interfaces

type BufferPool interface {
	Get() []byte
	Put([]byte)
}
