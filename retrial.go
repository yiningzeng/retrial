package main

import (
	"github.com/fsnotify/fsnotify"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func IniDefaultConfig() {
	viper.SetDefault("logger.name", "yiningzeng.log ")
	viper.SetDefault("logger.withRotationTime", 24)
	viper.SetDefault("logger.withRotationCount", 60)

	viper.SetDefault("Collect.enable", true)
	viper.SetDefault("Detect.enable", true)
	viper.SetDefault("Detect.incremental", true)

	viper.SetDefault("mr.mysqlHost", "127.0.0.1")
	viper.SetDefault("mr.mysqlPort", 3306)
	viper.SetDefault("mr.mysqlDatabase", "aoidatav4")
	viper.SetDefault("mr.mysqlUser", "root")
	viper.SetDefault("mr.mysqlPassword", "root")

	viper.SetConfigName("retrial") //  设置配置文件名 (不带后缀)
	viper.AddConfigPath("/workspace/appName/")   // 第一个搜索路径
	viper.AddConfigPath("/workspace/appName1")  // 可以多次调用添加路径
	viper.AddConfigPath(".")               // 比如添加当前目录
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // 搜索路径，并读取配置数据
	if err != nil {
		logger.WithField("config", err).Fatal("Fatal error config file")
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.WithField("config", e.Name).Fatal("配置文件被修改")
		err = viper.ReadInConfig()
		if err != nil {
			logger.WithField("config", e.Name).Fatal("重新读取配置文件失败" + err.Error())
		}
	})
	LoggerInit()
}

// 用于采集人工复判的结果数据
func Collect() {
	if viper.GetBool("debug") {
		DoCollectMingrui()
		// Do jobs without params
		_ = gocron.Every(10).Minute().Do(DoCollectMingrui)
	} else {
		_ = gocron.Every(1).Day().At(viper.GetString("collect.collectTime")).Do(DoCollectMingrui)
	}
	gocron.Start()
}

// 用于检测过滤明锐数据
func Detect() {
	logger.Println("我是检测室")
}

func main() {
	IniDefaultConfig()
	MysqlIni()
	//bord, err := DoBoardQuery(12)
	//res := DoFaultsQuery(1213232)
	//logger.Debug(res, err)
	//logger.Debug(bord.TestDate.GetPathString())
	//sBytes, er := afero.ReadFile(afero.NewOsFs(), "/home/baymin/go/src/test/outtxt/TESTTTT@/20201121/A/0000000004OD.dat")
	//if er != nil {
	//	return
	//} else {
	//	connect := string(sBytes) // 这里已经是明锐的所有文本信息了
	//	aa :=GetImageRectangle("707", connect, "C406.*?\n")
	//	logger.Info(aa)
	//}

	//aa := viper.Get("mr.dataIdCol")
	//aa = viper.Get("mr.dataIdCol.707")
	//logger.Println(aa)
	//_ = gocron.Every(1).Second().Do(func() {
	//	logger.Error("asdsdsd")
	//	logger.WithField("fuck", time.Now()).Debug("running")
	//})
	//gocron.Start()
	//fileNae := "/home/baymin/daily-work/go/src/retrial/AOIBin/test/0000000002ND.dat"
	//FileGetMingRuiDBBoardID(fileNae, afero.NewOsFs())
	//FileRead("/home/baymin/daily-work/go/src/retrial/AOIBin/test/0000000002ND.dat")

	if viper.GetBool("Collect.enable") {
		go Collect()
	}
	if viper.GetBool("Detect.enable") {
		go Detect()
	}
	go WatchDir(".")
	//Crop("/home/baymin/daily-work/go/src/retrial/AOIBin/pcbimage/TESTTTT@/2020-11/17/NO1/__20201117192436.png",
	//	"/home/baymin/daily-work/go/src/retrial/AOIBin/pcbimage/TESTTTT@/2020-11/17/NO1/__20201117192436.png.png",
	//	image.Rectangle{
	//		Min: image.Pt(16295, 3681),
	//		Max: image.Pt(16295+57, 3681+99),
	//	})//(16295, 3681, 57, 99))

	select{}
	//res := DoQuery(92111)
	//fmt.Println(res)
	//text := `Hello 世界！123 Go.dat`
	//// 查找连续的小写字母
	//reg := regexp.MustCompile(`^.*.dat`)
	//fmt.Printf("%q\n", reg.FindAllString(text, -1))
	//MingRuiDataCleaning()
	//FileTest()
	//FileWatch("aaa")
}

