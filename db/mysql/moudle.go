package mysql

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/mikeqiao/ant/log"
	"github.com/mikeqiao/ant/net/proto"
)

var DB *STypeMap

type SType struct {
	Name string
	Key  string
	Data map[string]string
}

func (s *SType) Init() {
	s.Data = make(map[string]string)
}

func (s *SType) Update(h bool, db *DBMysql) {
	table := new(proto.DataTableConfig)
	ts := new(Types)
	for k, v := range s.Data {
		t := new(TypeConf)
		t.Name = k
		t.Dtype = v
		if nil != t {
			ts.Ctypes = append(ts.Ctypes, *t)
		}
	}
	if nil != ts {
		d, err := json.Marshal(ts)
		if nil != err {
			log.Debug("err:%v", err)
		} else {
			table.SetData(string(d))
			if h {
				table.SetKeyName(s.Name)
				sql := table.DBUpdate()
				db.UpdateData(sql)
			} else {
				table.SetName(s.Name)
				sql := table.DBInsert()
				db.InsertData(sql)
			}

		}

	}

}

type STypeMap struct {
	DB       *DBMysql
	Typelist map[string]*SType
}

func (s *STypeMap) InitDB() {
	s.Typelist = make(map[string]*SType)
	s.DB = new(DBMysql)
	s.DB.InitDB()
}

func (s *STypeMap) OnClose() {
	s.DB.OnClose()
}

func (s *STypeMap) GetByName(name string) *SType {
	if v, ok := s.Typelist[name]; ok {
		return v
	}
	return nil
}

func (s *STypeMap) SetData(st *SType) bool {
	if nil != st {
		s.Typelist[st.Name] = st
		return true
	}
	return false
}

func (s *STypeMap) InitData() {
	h := s.DB.HasTable("DataTableConfig")
	if !h {
		tc := new(proto.DataTableConfig)
		CreateTable(tc, s.DB)
	} else {
		t := new(proto.DataTableConfig)
		sql, keys := t.DBSelect()
		tables := new(proto.ArrayTableConfig)
		e := s.DB.SelectAllData(sql, keys, tables)
		if nil != e {
			log.Debug("err:%v", e)
		}
		for _, v := range tables.Data {
			name := v.Name
			data := v.Data
			tys := new(Types)
			err := json.Unmarshal([]byte(data), tys)
			if nil != err {
				log.Debug("err:%v, data:%v", err, data)
			} else {
				st := new(SType)
				st.Init()
				st.Name = name
				for _, t := range tys.Ctypes {
					st.Data[t.Name] = t.Dtype
				}
				if nil != st {
					s.Typelist[st.Name] = st
				}
			}
		}
	}
}

func (s *STypeMap) Register(d interface{}) bool {
	t := reflect.TypeOf(d)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() == reflect.Ptr {
		return false
	}
	table := t.Elem().Name()
	//判断是否是需要注册的table
	if !CheckData(table) {
		return false
	}
	var ast *SType
	var have bool
	if v, ok := s.Typelist[table]; ok {
		ast = v
		have = true
	}
	if have && nil != ast {
		//判断是否修改表
		cmap := make(map[string]string)
		dmap := make(map[string]string)
		addcmap := make(map[string]string)
		adddmap := make(map[string]string)
		var uniques []string
		for k := 0; k < t.Elem().NumField(); k++ {
			name := t.Elem().Field(k).Name
			kind := t.Elem().Field(k).Type.Kind()
			if CheckColumn(name) {
				ty, def := CheckType(name, kind)
				if v, ok := ast.Data[name]; ok {
					if ty != v {
						cmap[name] = ty
						dmap[name] = def
						ok := CheckKey(name)
						if ok {
							uniques = append(uniques, name)
						}
						ast.Data[name] = ty
					}
				} else {
					ok := CheckKey(name)
					if ok {
						uniques = append(uniques, name)
					}

					addcmap[name] = ty
					adddmap[name] = def
					ast.Data[name] = ty
				}
			}
		}
		var res bool
		if len(cmap) > 0 {
			//修改表
			for k, v := range cmap {
				def, _ := dmap[k]
				sql := SqlChangeDType(table, k, v, def)
				res = s.DB.Exec(sql)
			}
		}
		if len(addcmap) > 0 {
			//修改表
			for k, v := range addcmap {
				def, _ := adddmap[k]
				sql := SqlAddColumn(table, k, v, def)
				res = s.DB.Exec(sql)
			}
		}
		for _, v := range uniques {
			sql := SqlSetUnique(table, v)
			res = s.DB.Exec(sql)
		}
		if res {
			ast.Update(have, s.DB)
		}
	} else {
		//创建表
		_, ast = CreateTable(d, s.DB)
		if nil != ast {
			ast.Update(false, s.DB)
		}
	}
	s.Typelist[table] = ast
	return true
}

func NewDBMoudle() *STypeMap {
	DB = new(STypeMap)
	DB.InitDB()
	return DB
}

func SqlChangeName(table, oldname, newname, dtype string) (sql string) {
	sql = "ALTER TABLE " + table + " CHANGE " + oldname + " " + newname + " " + dtype
	return
}

func SqlChangeDType(table, name, dtype, def string) (sql string) {
	if "text" == dtype {
		sql = "ALTER TABLE " + table + " MODIFY COLUMN " + name + " " + dtype + " NOT NULL"
	} else {
		sql = "ALTER TABLE " + table + " MODIFY COLUMN " + name + " " + dtype + " NOT NULL DEFAULT " + def
	}
	return
}

func SqlAddColumn(table, name, dtype, def string) (sql string) {
	if "text" == dtype {
		sql = "ALTER TABLE " + table + " ADD " + name + " " + dtype + " NOT NULL"
	} else {
		sql = "ALTER TABLE " + table + " ADD " + name + " " + dtype + " NOT NULL DEFAULT " + def
	}

	return
}

func SqlSetPrimaryKey(table, name string) (sql string) {
	sql = "ALTER TABLE " + table + " ADD PRIMARY KEY(" + name + ")"
	return
}

func SqlSetUnique(table, name string) (sql string) {
	sql = "ALTER TABLE " + table + " ADD UNIQUE (" + name + ")"
	return
}

func CheckType(k string, v reflect.Kind) (tstr, def string) {
	switch v {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		tstr = "int"
		def = "'0'"
	case reflect.Int64, reflect.Uint64:
		tstr = "bigint"
		def = "'0'"
	case reflect.String:
		ok := CheckKey(k)
		if ok {
			tstr = "varchar(64)"
		} else if CheckData(k) {
			tstr = "text"
		} else {
			tstr = "varchar(255)"
		}
		def = "''"
	case reflect.Float32:
		tstr = "float"
		def = "'0'"
	case reflect.Float64:
		tstr = "double"
		def = "'0'"
	case reflect.Slice:
		tstr = "varbinary(max)"
		def = "''"
	default:
	}
	return
}

func CheckKey(k string) bool {
	if strings.HasSuffix(k, "ID") {
		return true
	}
	return false
}

func CheckColumn(k string) bool {
	if strings.HasPrefix(k, "XXX_") {
		return false
	}
	return true
}

func CheckData(k string) bool {
	if strings.HasPrefix(k, "Data") {
		return true
	}
	return false
}

func CreateTable(d interface{}, db *DBMysql) (bool, *SType) {
	t := reflect.TypeOf(d)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() == reflect.Ptr {
		return false, nil
	}
	table := t.Elem().Name()

	var buffer bytes.Buffer
	buffer.WriteString("CREATE TABLE IF NOT EXISTS `")
	buffer.WriteString(table)
	buffer.WriteString("`( ")
	id := "`id` bigint AUTO_INCREMENT NOT NULL ,"
	buffer.WriteString(id)
	st := new(SType)
	st.Init()
	st.Name = table
	var uniques []string
	for k := 0; k < t.Elem().NumField(); k++ {

		name := t.Elem().Field(k).Name
		if CheckColumn(name) {
			kind := t.Elem().Field(k).Type.Kind()

			ty, def := CheckType(name, kind)
			buffer.WriteString("`")
			buffer.WriteString(name)
			buffer.WriteString("` ")
			buffer.WriteString(ty)
			if "text" == ty {
				buffer.WriteString(" NOT NULL ")
			} else {
				buffer.WriteString(" NOT NULL DEFAULT ")
				buffer.WriteString(def)
			}
			buffer.WriteString(",")

			st.Data[name] = ty

			ok := CheckKey(name)
			if ok {
				uniques = append(uniques, name)
			}
		}
	}
	buffer.WriteString("PRIMARY KEY ( `id` )")
	buffer.WriteString(") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
	sql := buffer.String()
	res := db.CreateTable(sql)
	if res {
		for _, v := range uniques {
			sql := SqlSetUnique(table, v)
			db.Exec(sql)
		}
	}

	return res, st
}
