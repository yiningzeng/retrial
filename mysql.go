package main

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type t_faults struct {
	DBBoardID   int `gorm:"Column:DBBoardID"`
	SubBoardID   int `gorm:"Column:SubBoardID"`
	ComponentID   int `gorm:"Column:ComponentID"`
	RectType   int `gorm:"Column:RectType"`
	ComponentName   string `gorm:"Column:ComponentName"`
	BarCode   string `gorm:"Column:BarCode"`
	Model   string `gorm:"Column:Model"`
	Code   string `gorm:"Column:Code"`
	Package   string `gorm:"Column:Package"`
	ReportResult   int `gorm:"Column:ReportResult"`
	ConfirmResult   int `gorm:"Column:ConfirmResult"`
	userConfirmResult   string `gorm:"Column:userConfirmResult"`
	XDiff   int `gorm:"Column:XDiff"`
	YDiff   int `gorm:"Column:YDiff"`
	ReportResultStr   string `gorm:"Column:ReportResultStr"`
	ConfirmResultStr   string `gorm:"Column:ConfirmResultStr"`
}

func DoQuery(dbBoardId int) []t_faults {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local",
		viper.Get("mr.mysqlUser"),
		viper.Get("mr.mysqlPassword"),
		viper.Get("mr.mysqlHost"),
		viper.Get("mr.mysqlPort"),
		viper.Get("mr.mysqlDatabase"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	var res []t_faults
	er := db.Model(&t_faults{}).Where(&t_faults{DBBoardID: dbBoardId}).Find(&res)
	if er != nil {
		fmt.Println(er)
	}
	return res
}

func CollectQuery(dbBoardId uint) []t_faults {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local",
		viper.Get("mr.mysqlUser"),
		viper.Get("mr.mysqlPassword"),
		viper.Get("mr.mysqlHost"),
		viper.Get("mr.mysqlPort"),
		viper.Get("mr.mysqlDatabase"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//sqlDB, err := db.DB()
	//// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	//sqlDB.SetMaxIdleConns(10)
	//
	//// SetMaxOpenConns 设置打开数据库连接的最大数量。
	//sqlDB.SetMaxOpenConns(100)
	//
	//// SetConnMaxLifetime 设置了连接可复用的最大时间。
	//sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		fmt.Println(err)
	}
	var res []t_faults
	er := db.Model(&t_faults{}).Where("DBBoardID > ?", dbBoardId).Find(&res)
	if er != nil {
		fmt.Println(er)
	}
	return res
}
