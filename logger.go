package main

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

func LoggerInit() {
	path := viper.GetString("logger.name")
	/* 日志轮转相关函数
	`WithLinkName` 为最新的日志建立软连接
	`WithRotationTime` 设置日志分割的时间，隔多久分割一次
	WithMaxAge 和 WithRotationCount二者只能设置一个
	 `WithMaxAge` 设置文件清理前的最长保存时间
	 `WithRotationCount` 设置文件清理前最多保存的个数
	*/
	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
	if viper.GetBool("debug") {
		log.SetFormatter(&log.TextFormatter{})
		//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
		//同时写文件和屏幕
		//fileAndStdoutWriter := io.MultiWriter([]io.Writer{writer, os.Stdout}...)
		log.SetOutput(os.Stdout)
		log.SetLevel(log.DebugLevel)
	} else {
		writer, _ := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			// rotatelogs.WithMaxAge(time.Duration(180)*time.),
			rotatelogs.WithRotationCount(60),
			rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		)
		log.SetOutput(writer)
		log.SetLevel(log.InfoLevel)
	}
}
