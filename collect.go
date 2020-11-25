package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Guitarbum722/align"
	"github.com/msterzhang/gpool"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gocv.io/x/gocv"
	"image"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)
// 1. 明锐数据采集循环遍历
func DoCollectMingrui() {
	pool := gpool.New(viper.GetInt("collect.threadNum"))
	fs := afero.NewOsFs()
	t1 := time.Now()
	//通过Walk来遍历目录下的所有子目录
	_ = filepath.Walk(viper.GetString("mr.outtxt"), func(file string, info os.FileInfo, err error) error {
		isDir, err := afero.IsDir(fs, file)
		if err != nil {
			logger.Error(err.Error())
		}
		_, fileName := filepath.Split(file)
		if !isDir && strings.HasSuffix(fileName, viper.GetString("mr.fileSuffix")) && !strings.Contains(viper.GetString("mr.excludeFiles"), fileName) {
			pool.Add(1)
			go FileRead(file, fs, pool)
		}
		return nil
	})
	pool.Wait()
	logger.WithFields(logger.Fields{"采集耗时": time.Now().Sub(t1)}).Info("数据采集完成")
}

// 2. 读取明锐数据文本信息
func FileRead(fileName string, fs afero.Fs, pool *gpool.Pool) {
	defer pool.Done()
	if fs == nil {
		fs = afero.NewOsFs()
	}
	isExist, _ := afero.Exists(fs, fileName+".done")
	// 是否增量采集和是否已经采集过
	if viper.GetBool("collect.incremental") && isExist {
		logger.WithFields(logger.Fields{"fileName": fileName}).Debug("已执行过")
		return
	}
	logger.WithFields(logger.Fields{"fileName": fileName}).Debug("开始采集")
	sBytes, er := afero.ReadFile(fs, fileName)
	if er != nil {
		logger.WithFields(logger.Fields{"fileName": fileName, "error": er.Error()}).Error("FileRead:43:打开文件失败")
		return
	} else {
		connect := string(sBytes) // 这里已经是明锐的所有文本信息了
		// 1. 获取明锐数据的记录ID
		DBBoardID, imgNameDateStr, DataVersion := FileGetMingRuiDBBoardID(&connect)
		changeTime, err := time.Parse("2006-01-02 15:04:05", strings.ReplaceAll(imgNameDateStr, "/", "-"))
		if err != nil {
			logger.Error(err.Error())
		}
		if DBBoardID == -1 { // 表明这个数据查询不到记录，说明维修站还未复判 等待下一次处理吧
			//_ = afero.WriteFile(fs, fileName + ".done", nil, 0755)
			return
		}
		// 2. 查询数据库相应的记录信息
		board, err := DoBoardQuery(DBBoardID)
		if err !=nil {
			logger.WithFields(logger.Fields{"DBBoardID": DBBoardID}).Warn("查询不到数据")
			return
		}
		dataRes := DoFaultsQuery(DBBoardID)
		if len(dataRes) == 0 {
			logger.WithFields(logger.Fields{"DBBoardID": DBBoardID}).Warn("查询不到数据")
			_ = afero.WriteFile(fs, fileName + ".done", nil, 0755)
			return
		}
		src := fmt.Sprintf("%s/%s/%s/%s/__%s.png",
			viper.GetString("mr.pcbimg"),
			board.ProjectName,
			changeTime.Format("2006-01/02"),
			board.MachineName,
			changeTime.Format("20060102150405"))
		isExist, _ = afero.Exists(fs, src)
		if !isExist {
			logger.WithField("src", src).Error("图片不存在")
			return
		}
		img := gocv.IMRead(src, gocv.IMReadColor)
		defer img.Close()
		for _, v := range dataRes {
			logger.WithFields(logger.Fields{"ComponentName": v.ComponentName, "DBBoardID": DBBoardID}).Info("通过DBBoardID查询的结果")
			// 3. 每个结果生成一个文件夹
			baseDir := viper.GetString("collect.powerAiAssetsSavePath") + "/" + board.ProjectName + "/" + v.ComponentName // 这里应该按照项目名下再分原件
			//
			// 4. 根据每个元件名称v.ComponentName 用正则表达式从connect获取相应的元件位置截取保存到上面创建的文件夹，按照(ok/ng)-DBBoardID-SubBoardID-ComponentID.jpg来存储
			// 这里保存的规则是 aoi判断的结果，和人工的结果对比，如果一致那就是保存下来，如果不一致以人工的为准。
			// 比如报了一个移位的缺陷，那么先在该元件下新建一个移位的目录，如果人工判断也是移位那么就保存为ng，如果人工判断是ok那么就在移位目录下保存为ok
			savePath := baseDir + "/" + strconv.Itoa(v.ReportResult) + "-" + v.ReportResultStr
			err = fs.MkdirAll(savePath, 0755)
			if err != nil {
				logger.WithField("dir", baseDir).Error("创建文件夹失败")
				continue
			}
			id := ""
			asset := PowerAiAsset{
				Asset: Asset{
					Format: "jpg",
					State:  2,
					Type:   1,
				},
				Regions: []string{},
				Version: "sort-5.1.1",
			}
			if v.ConfirmResult !=0 { //ng
				asset.Asset.Sort = "NG"
				id = fmt.Sprintf("NG_%d-%d-%d", v.DBBoardID, v.SubBoardID, v.ComponentID)
			} else { //ok
				asset.Asset.Sort = "OK"
				id = fmt.Sprintf("OK_%d-%d-%d", v.DBBoardID, v.SubBoardID, v.ComponentID)
			}
			imgSavePath := fmt.Sprintf("%s/%s.jpg", savePath, id)
			jsonSavePath := fmt.Sprintf("%s/%s.json", savePath, id)
			isExist, _ = afero.Exists(fs, imgSavePath)
			if !isExist { //只有在不存在的时候再裁剪
				rect := GetImageRectangle(DataVersion, &connect, v.ComponentName+".*?\n")
				Crop(&img, rect, imgSavePath)
				asset.Asset.Size = Size{Width: rect.Max.X - rect.Min.X, Height: rect.Max.Y - rect.Min.Y}
				asset.Asset.Id = id
				asset.Asset.Name = id + ".jpg"
				asset.Asset.Path = "file:" + imgSavePath
				jsonBytes, err := json.Marshal(asset)
				if err != nil {
					logger.Error("PowerAiAsset转json失败")
				}
				var out bytes.Buffer
				err = json.Indent(&out, jsonBytes, "", "\t")
				if err != nil {
					_ = afero.WriteFile(fs, jsonSavePath, jsonBytes, 0755)
				} else {
					_ = afero.WriteFile(fs, jsonSavePath, out.Bytes(), 0755)
				}
			}
			// 5. 最后再生成对应的powerai的json文件，按照(ok/ng)-DBBoardID-SubBoardID-ComponentID.json来存储
		}
		_ = afero.WriteFile(fs, fileName + ".done", nil, 0755)
		runtime.GC()
		//compile := regexp.MustCompile()
		//
		//submatch := compile.FindAllSubmatch(contents, -1)
		//logger.Trace(connect)
		//_ = afero.WriteFile(fs, "/home/baymin/daily-work/fff/aa.txt", []byte(strings.Split(connect, "\n")[0]), 0755)
	}
}

// 明锐数据清洗
// 主要是加载一行数据，根绝传进来的列来用,分隔，取出对应列的数据
func DataCleaningMingRui(content *string, col int) string {
	n := []int{col}
	input := strings.NewReader(*content)
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
func FileGetMingRuiDBBoardID(content *string) (int, string, string) {
	DBBoardID := -1
	version := ""
	imgFileDate := ""
	lines := strings.Split(*content, "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.ReplaceAll(lines[i], "\n", "")//以'\n'为结束符读入一行
		if i == 0 {
			version = DataCleaningMingRui(&line, viper.GetInt("mr.dataVersion.versionCol"))
		} else if i == viper.GetInt("mr.dataVersion." + version + ".idRow") - 1 { // 这里为了外面配置文件里行数和列数都是从第一行开始，所以要减一
			imgFileDate = DataCleaningMingRui(&line,  viper.GetInt("mr.dataVersion." + version + ".imgDateCol"))
			DBBoardID, _ = strconv.Atoi(DataCleaningMingRui(&line, viper.GetInt("mr.dataVersion." + version + ".idCol")))
			logger.WithFields(logger.Fields{"DataVersion": version, "DBBoardID": DBBoardID}).Info("数据库的编号")
			break
		}
		logger.WithField("line connect", line).Trace("trace")
	}
	return DBBoardID, imgFileDate, version
}

func GetImageRectangle(dataVersion string, content *string, reg string) image.Rectangle{
	re := regexp.MustCompile(reg)
	out := strings.Split(re.FindStringSubmatch(*content)[0], ",")
	x, _ := strconv.Atoi(out[viper.GetInt("mr.dataVersion."+dataVersion+".xCol") - 1])
	y, _ := strconv.Atoi(out[viper.GetInt("mr.dataVersion."+dataVersion+".yCol") - 1])
	w, _ := strconv.Atoi(out[viper.GetInt("mr.dataVersion."+dataVersion+".wCol") - 1])
	h, _ := strconv.Atoi( out[viper.GetInt("mr.dataVersion."+dataVersion+".hCol") - 1])
	logger.WithFields(logger.Fields{"x": x, "y": y, "w": w, "h": h}).Info("获取裁剪区域")
	return image.Rectangle{
			Min: image.Pt(x, y),
			Max: image.Pt(x+w, y+h)}
}