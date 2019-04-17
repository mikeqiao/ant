package mod

type Data interface {
	Close() bool
	Start() bool
	Run()
	GetKey() interface{}
}
