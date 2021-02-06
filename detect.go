package main

import (
	"github.com/fsnotify/fsnotify"
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Watch struct {
	watch *fsnotify.Watcher
}

//监控目录
func WatchDir(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
	}
	defer watcher.Close()
	//通过Walk来遍历目录下的所有子目录
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() && info.Name() != "log" {
			path, err := filepath.Abs(path)
			if err != nil {
				logger.Fatal(err)
			}
			err = watcher.Add(path)
			if err != nil {
				logger.Fatal(err)
			}
			//fmt.Println("监控 : ", path)
		}
		return nil
	})

	done := make(chan bool)
	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					logger.WithField("filename", ev.Name).Debug("Write")
					if strings.Contains(ev.Name, ".dat") {
						logger.WithField("filename", ev.Name).Debug("analysis")
						//
						//if err != nil {
						//	_ = errors.Wrap(err, "read failed")
						//}
						//defer conn.Close()
					}
				} else if ev.Op&fsnotify.Create == fsnotify.Create {
					logger.WithField("filename", ev.Name).Debug("Create")
					//获取新创建文件的信息，如果是目录，则加入监控中
					file, err := os.Stat(ev.Name)
					if err == nil && file.IsDir() {
						_ = watcher.Add(ev.Name)
						logger.WithField("filename", ev.Name).Debug("Add Watch")
					}
				} else if ev.Op&fsnotify.Remove == fsnotify.Remove {
					//如果删除文件是目录，则移除监控
					fi, err := os.Stat(ev.Name)
					if err == nil && fi.IsDir() {
						_ = watcher.Remove(ev.Name)
						logger.WithField("filename", ev.Name).Debug("Remove Watch")
					}
				} else if ev.Op&fsnotify.Rename == fsnotify.Rename {
					//如果重命名文件是目录，则移除监控 ,注意这里无法使用os.Stat来判断是否是目录了
					//因为重命名后，go已经无法找到原文件来获取信息了,所以简单粗爆直接remove
					logger.WithField("filename", ev.Name).Debug("Rename")
					_ = watcher.Remove(ev.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error(err)
			}
		}
	}()
	<-done
}

