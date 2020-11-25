package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type PowerAiAsset struct {
	Asset Asset `json:"asset"`
	Regions []string `json:"regions"`
	Version string `json:"version"`
}

type Asset struct {
	Format string `json:"format"`
	Id     string `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Size   Size   `json:"size"`
	State  int    `json:"state"`
	Type   int    `json:"type"`
	Sort   string `json:"sort"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type tBoards struct {
	DBBoardID         int       `gorm:"Column:DBBoardID"`
	BoardName         string    `gorm:"Column:BoardName"`
	ProjectName       string    `gorm:"Column:ProjectName"`
	MachineName       string    `gorm:"Column:MachineName"`
	TestDate          MyTime `gorm:"Column:TestDate"`
	TestBatch         string    `gorm:"Column:TestBatch"`
	SnInBatch         int       `gorm:"Column:SnInBatch"`
	TestTime          MyTime `gorm:"Column:TestTime"`
	TestCostTime      int       `gorm:"Column:TestCostTime"`
	ReportResult      int       `gorm:"Column:ReportResult"`
	ConfirmResult     int       `gorm:"Column:ConfirmResult"`
	ErrorComponent    int       `gorm:"Column:ErrorComponent"`
	TotalComponent    int       `gorm:"Column:TotalComponent"`
	ReportComponent   int       `gorm:"Column:ReportComponent"`
	AlarmComponent    int       `gorm:"Column:AlarmComponent"`
	WrongReportNum    int       `gorm:"Column:WrongReportNum"`
	ProgrammerName    string    `gorm:"Column:ProgrammerName"`
	TestOpName        string    `gorm:"Column:TestOpName"`
	ConfirmOpName     string    `gorm:"Column:ConfirmOpName"`
	SubBoardNum       int       `gorm:"Column:SubBoardNum"`
	ReportSubBoardNum int       `gorm:"Column:ReportSubBoardNum"`
	ErrSubBoardNum    int       `gorm:"Column:ErrSubBoardNum"`
	TrackName         string    `gorm:"Column:TrackName"`
}

type tFaults struct {
	DBBoardID         int    `gorm:"Column:DBBoardID"`
	SubBoardID        int    `gorm:"Column:SubBoardID"`
	ComponentID       int    `gorm:"Column:ComponentID"`
	RectType          int    `gorm:"Column:RectType"`
	ComponentName     string `gorm:"Column:ComponentName"`
	BarCode           string `gorm:"Column:BarCode"`
	Model             string `gorm:"Column:Model"`
	Code              string `gorm:"Column:Code"`
	Package           string `gorm:"Column:Package"`
	ReportResult      int    `gorm:"Column:ReportResult"`
	ConfirmResult     int    `gorm:"Column:ConfirmResult"`
	UserConfirmResult string `gorm:"Column:userConfirmResult"`
	XDiff             int    `gorm:"Column:XDiff"`
	YDiff             int    `gorm:"Column:YDiff"`
	ReportResultStr   string `gorm:"Column:ReportResultStr"`
	ConfirmResultStr  string `gorm:"Column:ConfirmResultStr"`
}

//MyTime 自定义时间
type MyTime time.Time
func (t *MyTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse("2006-01-02 15:04:05", timeStr)
	*t = MyTime(t1)
	return err
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

func (t MyTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format("2006-01-02 15:04:05"), nil
}

func (t *MyTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = MyTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *MyTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

func (t *MyTime) GetPathString() string {
	return fmt.Sprintf("%s", time.Time(*t).Format("2006-01/02"))
}

func (t *MyTime) GetImageNameString() string {
	return fmt.Sprintf("__%s.png", time.Time(*t).Format("20060102150405"))
}