package mysql

import (
	"fmt"

	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net/proto"
	"github.com/mikeqiao/ant/rpc"
)

func HandleDBRQ(call *rpc.CallInfo) {

	in := call.Args.(*proto.DBServerRQ) //本服务需要的参数
	//		call.Data //附加信息
	t := in.GetType()
	sql := in.GetSql()
	keys := in.GetKey()
	msgid := in.GetMsgid()

	data := call.Data
	var res interface{}
	if nil != data && nil != data.Agent && nil != data.Agent.Processor {
		res = data.Agent.Processor.GetMsg(uint16(msgid))
	}

	var ok bool
	var err error
	if nil != DB && nil != DB.DB {
		switch t {
		case DB_insert:
			ok = DB.DB.InsertData(sql)
		case DB_del:
			ok = DB.DB.DeleteData(sql)
		case DB_update:
			ok = DB.DB.UpdateData(sql)
		case DB_select:
			err = DB.DB.SelectData(sql, keys, res)
		case DB_select_all:
			err = DB.DB.SelectAllData(sql, keys, res)
		}
	} else {
		err = fmt.Errorf("DB is nil")
	}
	if msgid > 0 {
		cb := new(proto.DBServerRS)
		if ok {
			cb.Result = 0
		} else {
			cb.Result = 1
		}
		call.SetResult(res, nil)
	} else {
		call.SetResult(res, err)
	}
	//	计算得出结果
	log.Debug("id:%v, ok:%v, err:%v", msgid, ok, err)

}
