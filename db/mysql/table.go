package mysql

type DATA interface {
	GetFieldID(ks []string) (p []interface{})
}

type Types struct {
	Ctypes []TypeConf
}

type TypeConf struct {
	Name  string
	Dtype string
}
