package conf

import (
	"time"
)

var Config struct {
	LenStackBuf int

	// log
	LogLevel int
	LogPath  string
	LogFlag  int

	// netClient
	ConnConf []ClientConfig

	// netServer
	ServerConf []ServerConfig

	ConnectInterval time.Duration

	//websocket
	HTTPTimeout time.Duration

	LittleEndian bool

	ServiceInfo ServiceConfig

	Dbconfig DBConfig

	Redisconfig RedisConfig
}

type DBConfig struct {
	UserName string
	Password string
	Ip       string
	Port     string
	DbName   string
}

type RedisConfig struct {
	Ip       string
	Port     string
	Password string
	MaxIdle  int32
	Life     int32 //day
}

type ServiceConfig struct {
	ServerId      int64
	ServerType    int32
	ServerVersion int32
	ServerName    string
	GameType      int32
}

type ClientConfig struct {
	UId             int64
	Name            string
	Nettype         int //网络通信类型 1 tcp        2 websocket
	PendingWriteNum int
	MaxMsgLen       uint32
	Processortype   int //消息编码类型 1 protobuf        2 json
	ConnAddr        string
	// tcp
	LenMsgLen int   //消息长度字段占用字节数
	Pid       int32 //processor ID
}

type ServerConfig struct {
	UId             int64
	Name            string
	Nettype         int //网络通信类型 1 tcp        2 websocket
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processortype   int    //消息编码类型 1 protobuf        2 json
	ListenAddr      string //监听地址
	LenMsgLen       int    //消息长度字段占用字节数
	Pid             int32  //processor ID
}

func SetServerConf(s_conf ...ServerConfig) {
	for j := 0; j < len(s_conf); j++ {
		Config.ServerConf = append(Config.ServerConf, s_conf[j])
	}
}
