package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type t_faults struct {
	DBBoardID   uint
	SubBoardID   uint
	ComponentID   uint
	RectType   uint
	ComponentName   string
	BarCode   string
	Model   string
	Code   string
	Package   string
	ReportResult   uint
	ConfirmResult   uint
	userConfirmResult   string
	XDiff   uint
	YDiff   uint
	ReportResultStr   string
	ConfirmResultStr   string
}

func doQuery(dbBoardId uint) []t_faults {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:root@tcp(192.168.31.77:3306)/aoidatav4?parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	res := []t_faults{}

	er := db.Find(&res).Where(t_faults{DBBoardID: dbBoardId})
	if er != nil {
		fmt.Println(er)
	}
	return res
}

