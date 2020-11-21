package main

import (
	"bufio"
	"bytes"
	"github.com/Guitarbum722/align"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 1. 明锐数据采集循环遍历
func DoCollectMingrui() {
	//通过Walk来遍历目录下的所有子目录
	_ = filepath.Walk(viper.GetString("mr.outtxt"), func(file string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if !info.IsDir() && filepath.Ext(info.Name()) == ".dat" {
			//path, fileName := filepath.Split(file)
			//logger.WithFields(logger.Fields{"path": path, "fileName": fileName}).Debug("start .dat")
			//go FileReadLine(file)
			go FileRead(file)
		}
		return nil
	})
}

// 2. 读取明锐数据文本信息
func FileRead(fileName string) {
	fs := afero.NewOsFs()
	isExist, _ := afero.Exists(fs, fileName+".done")
	// 是否增量采集和是否已经采集过
	if viper.GetBool("collect.incremental") && isExist {
		logger.WithFields(logger.Fields{"fileName": fileName}).Debug("已执行过")
		return
	}
	logger.WithFields(logger.Fields{"fileName": fileName}).Debug("开始采集")
	sBytes, er := afero.ReadFile(fs, fileName)
	if er != nil {
		logger.WithField("fileName", fileName).Error("打开文件失败")
		return
	} else {
		connect := string(sBytes) // 这里已经是明锐的所有文本信息了
		// 1. 获取明锐数据的记录ID
		DBBoardID := FileGetMingRuiDBBoardID(fileName, fs)
		if DBBoardID == -1 { // 表明这个数据查询不到记录，就直接标记为已处理
			_ = afero.WriteFile(fs, fileName + ".done", nil, 0755)
			return
		}
		// 2. 查询数据库相应的记录信息
		for _, v := range DoQuery(DBBoardID) {
			logger.WithFields(logger.Fields{"ComponentName": v.ComponentName, "DBBoardID": DBBoardID}).Info("通过DBBoardID查询的结果")
			// 3. 每个结果生成一个文件夹
			dir := viper.GetString("collect.powerAiAssetsSavePath") + "/" +  v.ComponentName  这里应该按照项目名下再分原件
			err := fs.MkdirAll(dir, 0755)
			if err != nil {
				logger.WithField("dir", dir).Error("创建文件夹失败")
				continue
			}
			// 4. 根据每个元件名称v.ComponentName 用正则表达式从connect获取相应的元件位置截取保存到上面创建的文件夹，按照(ok/ng)-DBBoardID-SubBoardID-ComponentID.jpg来存储

			// 5. 最后再生成对应的powerai的json文件，按照(ok/ng)-DBBoardID-SubBoardID-ComponentID.json来存储

		}

		//compile := regexp.MustCompile()
		//
		//submatch := compile.FindAllSubmatch(contents, -1)
		logger.Trace(connect)
		//_ = afero.WriteFile(fs, "/home/baymin/daily-work/go/src/retrial/AOIBin/test/" + onlyName, sBytes, 0755)
	}
}

// 明锐数据清洗
// 主要是加载一行数据，根绝传进来的列来用,分隔，取出对应列的数据
func DataCleaningMingRui(content string, col int) string {
	n := []int{col}
	input := strings.NewReader(content)
	output := bytes.NewBufferString("")
	bb := output.String()
	logger.Println(bb)
	aligner := align.NewAlign(input, output, ",", align.TextQualifier{}) // io.Reader, io.Writer, input delimiter, text qualifier
	//aligner.OutputSep("|")
	aligner.FilterColumns(n)
	aligner.UpdatePadding(
		align.PaddingOpts{
			Justification:  align.JustifyLeft,
			Pad: 0, // default is 1
		},
	)
	aligner.Align()
	return strings.ReplaceAll(output.String(), "\n", "")
}

// 通过传入的.dat文件名获取数据库对应的DBBoardID
func FileGetMingRuiDBBoardID(fileName string, fs afero.Fs) int {
	DBBoardID := -1
	file, err := fs.Open(fileName)
	if err != nil {
		logger.WithField("fileName", fileName).Error("打开文件失败")
		return DBBoardID
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	readRow := 1
	version := ""
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		if readRow == 1 {
			version = DataCleaningMingRui(line, viper.GetInt("mr.dataVersion.versionCol"))
		} else if readRow == viper.GetInt("mr.dataVersion." + version + ".idRow") {
			DBBoardID, _ = strconv.Atoi(DataCleaningMingRui(line, viper.GetInt("mr.dataVersion." + version + ".idCol")))
			logger.WithFields(logger.Fields{"DataVersion": version, "DBBoardID": DBBoardID, "FileName": fileName}).Info("数据库的编号")
			break
		}
		readRow++
		logger.WithField("line connect", line).Trace("trace")
	}
	return DBBoardID
}

func FileReadLine(fileName string) {
	fs := afero.NewOsFs()
	isExist, _ := afero.Exists(fs, fileName+".done")
	// 是否增量采集和是否已经采集过
	if viper.GetBool("collect.incremental") && isExist {
		logger.WithFields(logger.Fields{"fileName": fileName}).Debug("已执行过")
		return
	}
	logger.WithFields(logger.Fields{"fileName": fileName}).Debug("开始采集")
	sBytes, er := afero.ReadFile(fs, fileName)
	if er != nil {
		logger.WithField("fileName", fileName).Error("打开文件失败")
		return
	} else {
		_, onlyName := filepath.Split(fileName)
		_ = afero.WriteFile(fs, "/home/baymin/daily-work/go/src/retrial/AOIBin/test/" + onlyName, sBytes, 0755)
	}
	file, err := fs.Open(fileName)
	if err != nil {
		logger.WithField("fileName", fileName).Error("打开文件失败")
		return
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			_ = afero.WriteFile(fs, fileName+".done", nil, 0755)
			break
		}
		logger.WithField("line connect", line).Trace("trace")
	}
}
