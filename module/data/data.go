package mod

type Data interface {
	Close() bool
	Start() bool
	GetKey() interface{}
}
