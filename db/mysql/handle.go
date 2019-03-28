package mysql

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/log"
)

const (
	DB_action int32 = iota
	DB_add
	DB_del
	DB_update
	DB_find
	DB_inc
	DB_find_first
	DB_find_or_create
)

func HandlDBAction(rtype int32, column string, key, value interface{}) int32 {
	log.Debug("%v, %v", rtype, column)
	if nil == DB {
		return 1
	}
	result := int32(1)
	res := true
	if 1 == rtype {
		res = DB.AddData(value)
	} else if 2 == rtype {
		res = DB.DelData(key, value)
	} else if 3 == rtype {
		res = DB.UpdataData(key, value)
	} else if 4 == rtype {
		res = DB.FindData(key, value)
	} else if 5 == rtype {
		res = DB.IncData(column, value)
	} else if 6 == rtype {
		res = DB.FindDataOnly(value)
	} else if 7 == rtype {
		res = DB.FindOrCreateData(key, value)
	} else if 8 == rtype {
		res = DB.FindDataAll(key, value)
	} else {
		log.Debug("HandleServerDBMsgReq type error: %v", rtype)
		res = false
	}
	if res {
		result = 0
	}
	return result

}

func DbConn() *gorm.DB {
	UserName := conf.Config.Dbconfig.UserName
	Password := conf.Config.Dbconfig.Password
	Ip := conf.Config.Dbconfig.Ip
	Port := conf.Config.Dbconfig.Port
	DbName := conf.Config.Dbconfig.DbName
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", UserName, Password, Ip, Port, DbName)
	db, err := gorm.Open("mysql", connArgs)

	if err != nil {
		log.Fatal("err:%v", err)
	}
	db.SingularTable(true)
	return db
}

func CreateKey(key string) string {
	//大写字符的数值是 'A'=65到 ''Z=90
	nkey := ""
	for k, v := range key {
		if v >= 65 && v <= 90 {
			if k > 0 && !(v == 'D' && key[k-1] == 'I') && !(v == 'I' && key[k-1] == 'U') {
				nkey += "_"
			}
			nkey += string(v + 32)
		} else {
			nkey += string(v)
		}
	}
	return nkey

}

func CheckData(key string) bool {
	//大写字符的数值是 'A'=65到 ''Z=90
	res := strings.Index(key, "Data")
	if 0 == res {
		return true
	}
	return false

}

func GetKeyColumn(key string) (string, int32) {
	k := CreateKey(key)
	l := len(key)
	if l > 2 && 'D' == key[l-1] && 'I' == key[l-2] && 'U' == key[l-3] {
		return k, 1
	} else if l > 1 && 'D' == key[l-1] && 'I' == key[l-2] {
		return k, 2
	}
	return "", 0
}
