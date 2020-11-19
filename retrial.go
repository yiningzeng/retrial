package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"github.com/Guitarbum722/align"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rjeczalik/notify"
	"github.com/spf13/afero"
	"io"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func MingRuiDataCleaning() {
	n := []int{3,1,100}
	input := strings.NewReader("first,middle,last, 123123, 343,5,65,7,6,8,9,8,09,0,9,,6,324,21,321,3,213,12,32,,,,,,12323,")
	output := bytes.NewBufferString("")
	aligner := align.NewAlign(input, output, ",", align.TextQualifier{
		On:        true,
		Qualifier: "\"",
	}) // io.Reader, io.Writer, input delimiter, text qualifier
	//aligner.OutputSep("|")
	aligner.FilterColumns(n)
	aligner.Align()
	//aligner.Align()
	aa := output.String()
	log.Println("Got event:", aa)
}

func FileTest() {
	fs := afero.NewOsFs()
	//fs := afero.NewRegexpFs(afero.NewOsFs(), regexp.MustCompile(`^.*.dat`))
	//_, err := fs.Create("file.dat")
	list, err := afero.ReadDir(fs, "/home/baymin/go/src/test/outtxt/TESTTTT@/20201116/A/")
	for i :=0; i< len(list); i++ {
		if list[i].IsDir() {
			continue
		}
		fmt.Println(list[i].Name())
	}

	file, err := fs.Open("/home/baymin/go/src/test/outtxt/TESTTTT@/20201116/A/0000000010NY.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		log.Println(line)
	}
	file.Close()
	//aa := b.Buffer.String()
}

func FileWatch(saveDir string) {
	fs := afero.NewOsFs()
	fs.Name()
	//定义变量，用于接收命令行的参数值
	var dir string
	//&user 就是接收用户命令行中输入的 -u 后面的参数值
	//"u" ,就是 -u 指定参数
	//"" , 默认值
	//"用户名,默认为空" 说明
	flag.StringVar(&dir, "dir", "./...", "监听的文件夹")
	//转换
	flag.Parse()
	// Make the channel buffered to ensure no event is dropped. Notify will drop
	// an event if the receiver is not able to keep up the sending pace.
	c := make(chan notify.EventInfo, 1)

	// Set up a watchpoint listening on events within current working directory.
	// Dispatch each create and remove events separately to c.
	for {
		if err := notify.Watch(dir, c, notify.All); err != nil {
			log.Fatal(err)
		}
		ei := <-c
		if ei.Event() == notify.Create {
			log.Println("新增 Got event:", ei)
			_, fileName := filepath.Split(ei.Path())
			fmt.Println(fileName)
		} else if ei.Event() == notify.Remove {
			log.Println("删除 Got event:", ei)
		} else if ei.Event() == notify.Rename {
			log.Println("重命名 Got event:", ei)
		} else if ei.Event() == notify.Write {
			log.Println("写入 Got event:", ei)
		}
	}
}

func MysqlConnect() {
	db, err := sql.Open("mysql", "root:root@tcp(192.168.31.77:3306)/aoidatav4")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	users := sq.Select("*").From("t_faults")

	//active := users.Where(sq.Eq{"DBBoardID": 92108, "Model": "C0603-棕色"})
	sqlStr, args, err := users.Where(sq.Eq{"DBBoardID": 92108}).ToSql()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// Execute the query
	rows, err := db.Query(sqlStr, args...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
func main() {
	res := doQuery(92111)
	fmt.Println(res)
	MysqlConnect()
	text := `Hello 世界！123 Go.dat`
	// 查找连续的小写字母
	reg := regexp.MustCompile(`^.*.dat`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	MingRuiDataCleaning()
	FileTest()
	FileWatch("aaa")
}
