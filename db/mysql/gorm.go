package mysql

import (
	"encoding/json"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/mikeqiao/ant/log"
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

func (s *SType) Update(db *gorm.DB) {
	table := new(TableConfig)
	table.Name = s.Name
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
			table.Data = d

			db.Where(TableConfig{Name: s.Name}).Assign(TableConfig{Data: d}).FirstOrCreate(table)
		}

	}

}

type STypeMap struct {
	DB       *gorm.DB
	Typelist map[string]*SType
}

func (s *STypeMap) InitDB() {
	s.DB = DbConn()
	s.Typelist = make(map[string]*SType)
}

func (s *STypeMap) OnClose() {
	s.DB.Close()
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
	h := s.DB.HasTable(&TableConfig{})
	if !h {
		e := s.DB.AutoMigrate(&TableConfig{}).Error
		if nil != e {
			log.Debug("err:%v", e)
		}
	} else {
		var tables []TableConfig
		e := s.DB.Find(&tables).Error
		if nil != e {
			log.Debug("err:%v", e)
		}
		for _, v := range tables {
			name := v.Name
			data := v.Data
			tys := new(Types)
			err := json.Unmarshal(data, tys)
			if nil != err {
				log.Debug("err:%v, data:%v", e, tys)
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

func (s *STypeMap) Register(data interface{}) bool {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() == reflect.Ptr {
		return false
	}
	//判断是否是需要注册的table
	if !CheckData(t.Elem().Name()) {
		return false
	}

	st := new(SType)
	st.Data = make(map[string]string)
	st.Name = t.Elem().Name()
	var ast *SType
	if v, ok := s.Typelist[t.Elem().Name()]; ok {
		ast = v
	}
	//比较不同 记录 区别 执行修改
	Mmap := make(map[string]string)
	var auto = false
	var keys string
	var uniques []string
	kmap := make(map[string]string)
	tablename := CreateKey(t.Elem().Name())
	//	keys = append(keys, tablename)
	for k := 0; k < t.Elem().NumField(); k++ {
		//	fmt.Printf("%v, %s \n", t.Field(k).Type, t.Field(k).Name)
		st.Data[t.Elem().Field(k).Name] = t.Elem().Field(k).Type.Name()
		if nil != ast {
			if tt, ok := ast.Data[t.Elem().Field(k).Name]; ok {
				if tt != t.Elem().Field(k).Type.Name() {
					// ModifyColumn 表
					Mmap[t.Elem().Field(k).Name] = t.Elem().Field(k).Type.Name()
					//		e := db.Model(row).ModifyColumn("server_id", "int").Error
				}
			} else {
				//AutoMigrate表
				auto = true
				key, res := GetKeyColumn(t.Elem().Field(k).Name)
				if 1 == res {
					if "" == keys {
						keys += key
					} else {
						keys += ","
						keys += key
					}
					kmap[t.Elem().Field(k).Name] = t.Elem().Field(k).Type.Name()

				} else if 2 == res {
					uniques = append(uniques, key)
					kmap[t.Elem().Field(k).Name] = t.Elem().Field(k).Type.Name()
				}
			}
		} else {
			//create表
			auto = true
			key, res := GetKeyColumn(t.Elem().Field(k).Name)
			if 1 == res {
				if "" == keys {
					keys += key
				} else {
					keys += ","
					keys += key
				}
				kmap[t.Elem().Field(k).Name] = t.Elem().Field(k).Type.Name()
			} else if 2 == res {
				uniques = append(uniques, key)
				kmap[t.Elem().Field(k).Name] = t.Elem().Field(k).Type.Name()
			}
		}

	}

	if auto {

		e := s.DB.AutoMigrate(data).Error
		if nil != e {
			log.Debug("err:%v", e)
		}
		for k, v := range kmap {

			if "string" == v {
				tstr := "varchar(64)"
				nk := CreateKey(k)
				e := s.DB.Model(data).ModifyColumn(nk, tstr).Error
				if nil != e {
					log.Debug("err:%v", e)
				}
			}

		}

		if "" != keys {
			//		alter table table_test add primary key(id);
			sqlstr := "alter table " + tablename + " add primary key(" + keys + ")"
			err := s.DB.Exec(sqlstr).Error
			if nil != err {
				log.Debug("err:%v, sql:%v", err, sqlstr)
			}
		}
		for _, v := range uniques {
			//		alter table table_test add primary key(id);
			sqlstr := "alter table " + tablename + " add unique(" + v + ")"
			err := s.DB.Exec(sqlstr).Error
			if nil != err {
				log.Debug("err:%v, sql:%v", err, sqlstr)
			}
		}

	}

	for k, v := range Mmap {
		tstr := ""

		switch v {
		case "int":
			tstr = "int"
			break
		case "int64":
			tstr = "bigint"
			break
		case "string":
			tstr = "varchar"
			break
		case "float64":
			tstr = "double"
			break
		}
		if "" != tstr {
			nk := CreateKey(k)
			e := s.DB.Model(data).ModifyColumn(nk, tstr).Error
			if nil != e {
				log.Debug("err:%v", e)
			}
		}

	}
	//记录新的内容
	st.Update(s.DB)
	s.Typelist[t.Elem().Name()] = st
	return true
}

func (s *STypeMap) AddData(args interface{}) bool {
	e := s.DB.Create(args).Error
	if nil != e {
		log.Debug("err:%v", e)
		return false
	}
	return true
}

func (s *STypeMap) DelData(where, args interface{}) bool {
	e := s.DB.Model(args).Where(where).Delete(args).Error
	if nil != e {
		log.Debug("err:%v", e)
		return false
	}
	return true
}

func (s *STypeMap) UpdataData(where, args interface{}) bool {
	e := s.DB.Model(args).Where(where).Updates(args).Error
	if nil != e {
		log.Debug("err:%v", e)
		return false
	}
	return true
}

func (s *STypeMap) FindData(where, args interface{}) bool {
	e := s.DB.Model(args).Where(where).First(args).Error
	if nil != e {
		log.Debug("err:%v, key:%v, value:%v", e, where, args)
		return false
	}
	return true
}

func (s *STypeMap) FindOrCreateData(where, args interface{}) bool {
	e := s.DB.Model(args).Where(where).FirstOrCreate(args).Error
	if nil != e {
		log.Debug("err:%v, key:%v, value:%v", e, where, args)
		return false
	}
	return true
}

func (s *STypeMap) FindDataOnly(args interface{}) bool {
	e := s.DB.Model(args).First(args).Error
	if nil != e {
		log.Debug("err:%v, key:%v, value:%v", e, args)
		return false
	}
	return true
}

func (s *STypeMap) IncData(column string, args interface{}) bool {
	e := s.DB.Model(args).UpdateColumn(column, gorm.Expr(column+" + ?", 1)).First(args).Error
	if nil != e {
		log.Debug("err:%v, key:%v, value:%v", e, column, args)
		return false
	}
	return true
}

//传递的args 必须是个 []struct类型
func (s *STypeMap) FindDataAll(where, args interface{}) bool {
	t := reflect.TypeOf(args)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() == reflect.Ptr {
		return false
	}

	if t.Elem().NumField() <= 0 {
		return false
	}
	res := reflect.New(t.Elem().Field(0).Type).Interface()
	e := s.DB.Model(where).Where(where).Find(res).Error
	log.Debug("res:%v", res)

	if nil != e {
		log.Debug("err:%v, key:%v, value:%v", e, where, res)
		return false
	}
	v := reflect.ValueOf(args).Elem()
	dst := v.FieldByName(t.Elem().Field(0).Name)
	if reflect.TypeOf(res).Elem() != dst.Type() || !dst.CanSet() {
		log.Debug("err type")
		return false
	}
	dst.Set(reflect.ValueOf(res).Elem())

	return true
}
