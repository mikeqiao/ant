package msgtype

const (
	_None uint16 = iota
	NewConnect
	DelConnect
	ServerDelConnect
	ServerTick
	ServerLoginRQ
	ServerLoginRS
	ServerRegister
	ServerDelFunc
	ServerCall
	ServerCallBack
)
