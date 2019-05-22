package mysql

import (
	"fmt"
	"reflect"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	conf "github.com/mikeqiao/ant/config"
	"github.com/mikeqiao/ant/log"
)

const (
	DB_action int32 = iota
	DB_insert
	DB_del
	DB_update
	DB_select
	DB_select_all
	DB_inc
)

//Db数据库连接池
type DBMysql struct {
	DB *sql.DB
}

func (d *DBMysql) InitDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{conf.Config.Dbconfig.UserName, ":", conf.Config.Dbconfig.Password, "@tcp(", conf.Config.Dbconfig.Ip, ":", conf.Config.Dbconfig.Port, ")/", conf.Config.Dbconfig.DbName, "?charset=utf8"}, "")

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, terr := sql.Open("mysql", path)
	if terr != nil {
		log.Debug("err:%v", terr)
		return
	}
	d.DB = DB
	//设置数据库最大连接数
	d.DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	d.DB.SetMaxIdleConns(10)
	//验证连接
	if err := d.DB.Ping(); err != nil {
		log.Debug("opon database fail")
		return
	}
	log.Debug("connnect mysql success")
}

func (d *DBMysql) OnClose() {
	d.DB.Close()
}

func (d *DBMysql) Exec(sqlstr string) bool {
	_, err := d.DB.Exec(sqlstr)
	if err != nil {
		log.Error("Exec fail sql:%v", sqlstr)
		return false
	}
	return true
}

func (d *DBMysql) InsertData(sqlstr string) int64 {
	res, err := d.DB.Exec(sqlstr)
	if err != nil {
		log.Error("Exec fail sql:%v", sqlstr)
		return 0
	}
	id, _ := res.LastInsertId()
	return id
}

func (d *DBMysql) DeleteData(sqlstr string) int64 {
	res, err := d.DB.Exec(sqlstr)
	if err != nil {
		log.Error("Exec fail sql:%v", sqlstr)
		return 0
	}
	count, _ := res.RowsAffected()
	return count
}

func (d *DBMysql) UpdateData(sqlstr string) int64 {
	res, err := d.DB.Exec(sqlstr)
	if err != nil {
		log.Error("Exec fail sql:%v", sqlstr)
		return 0
	}
	count, _ := res.RowsAffected()
	return count
}

func (d *DBMysql) SelectData(sqlstr string, k []string, res interface{}) error {
	if nil == res {
		log.Error("error type, res is nil")
		err := fmt.Errorf("error type, res is nil")
		return err
	}
	a, ok := res.(DATA)
	if !ok {
		log.Error("SelectData res:%v, sql:%v", res, sqlstr)
		err := fmt.Errorf("error type")
		return err
	}
	address := a.GetFieldID(k)
	err := d.DB.QueryRow(sqlstr).Scan(address...)
	if err != nil {
		log.Error("SelectData err:%v, sql:%v", err, sqlstr)
	}
	return err
}

func (d *DBMysql) SelectAllData(sqlstr string, b []string, res interface{}) (err error) {
	//执行查询语句
	if nil == res {
		log.Error("error type, res is nil")
		err = fmt.Errorf("error type, res is nil")
		return
	}
	rows, err := d.DB.Query(sqlstr)
	if err != nil {
		log.Error("SelectAllData err:%v, sql:%v", err, sqlstr)
		err = fmt.Errorf("SelectAllData err:%v, sql:%v", err, sqlstr)
		return
	}
	//循环读取结果
	t := reflect.TypeOf(res)
	t1 := t.Elem().Field(0)
	t2 := t1.Type.Elem().Elem()
	name := t1.Name
	v := reflect.ValueOf(res).Elem()
	v1 := v.Field(0)
	v2 := v.FieldByName(name)
	for rows.Next() {
		info := reflect.New(t2).Interface()
		a, ok := info.(DATA)
		if !ok {
			log.Error("SelectAllData res:%v, sql:%v", info, sqlstr)
			err = fmt.Errorf("error type")
			return
		}
		add := a.GetFieldID(b)
		rows.Scan(add...)
		//		res = append(res, info)
		mVal := reflect.ValueOf(info)
		v1 = reflect.Append(v1, mVal)
	}
	v2.Set(v1)

	if err != nil {
		return
	}
	return

}

func (d *DBMysql) HasTable(table string) bool {
	//开启事务
	sqlstr := "SELECT count(*) FROM INFORMATION_SCHEMA.tables WHERE table_name = '" + table + "'"
	var count int32
	err := d.DB.QueryRow(sqlstr).Scan(&count)
	if err != nil {
		log.Error("HasTable err:%v, sql:%v", err, sqlstr)
	}
	if count > 0 {
		return true
	}
	return false
}

func (d *DBMysql) CreateTable(sqlstr string) bool {
	_, err := d.DB.Exec(sqlstr)
	if err != nil {
		log.Error("Exec fail sql:%v", sqlstr)
		return false
	}
	return true
}
