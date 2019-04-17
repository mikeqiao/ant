package net

type Processor interface {
	//发送处理
	Route(id uint16, msg interface{}, userData *UserData) error
	//解包数据
	Unmarshal(data []byte) (uint16, interface{}, error)
	//打包数据
	Marshal(msg interface{}) (uint16, [][]byte, error)

	Register(msg interface{}, id uint16) uint16

	GetMsg(id uint16) interface{}
}
