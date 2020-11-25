package main

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func MysqlIni() bool {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local",
		viper.Get("mr.mysqlUser"),
		viper.Get("mr.mysqlPassword"),
		viper.Get("mr.mysqlHost"),
		viper.Get("mr.mysqlPort"),
		viper.Get("mr.mysqlDatabase"))
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.WithFields(logger.Fields{"err": "gorm.open"}).Error(err.Error())
		return false
	}
	return true
}

func DoBoardQuery(dbBoardId int) (t tBoards, err error) {
	var tBs tBoards
	result := db.Model(&tBoards{}).Where(&tBoards{DBBoardID: dbBoardId}).First(&tBs)
	if result.Error != nil {
		logger.WithFields(logger.Fields{"DBBoardID": dbBoardId}).Warn(result.Error)
	}
	return tBs, result.Error
}

func DoFaultsQuery(dbBoardId int) []tFaults {
	var tFs []tFaults
	result := db.Model(&tFaults{}).Where(&tFaults{DBBoardID: dbBoardId}).Find(&tFs)
	if result.Error != nil {
		logger.WithFields(logger.Fields{"DBBoardID": dbBoardId}).Warn(result.Error)
	}
	return tFs
}

func CollectQuery(dbBoardId uint) []tFaults {
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
		logger.Error(err)
	}
	var res []tFaults
	er := db.Model(&tFaults{}).Where("DBBoardID > ?", dbBoardId).Find(&res)
	if er != nil {
		logger.Error(er)
	}
	return res
}
