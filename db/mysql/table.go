package mysql

import (
	"github.com/jinzhu/gorm"
)

type TableConfig struct {
	gorm.Model
	Name string `gorm:"primary_key;size:75"`
	Data []byte `gorm:"size:10000"`
}

type Types struct {
	Ctypes []TypeConf
}

type TypeConf struct {
	Name  string
	Dtype string
}
