package rdbms

import (
	"boot"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Mysql gorm.DB

func init() {
	name, conf := boot.App.DsNameConf("rdbms")
	var err error
	Mysql, err = gorm.Open(name, conf)
	if err != nil {
		fmt.Println(err)
	}
	Mysql.LogMode(true)
	if err != nil {
		panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
	}
}
