package serv

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/TruthHun/gpm/conf"
	"github.com/TruthHun/gpm/utils"
)

var (
	files sync.Map
	ext   sync.Map
	pid   []string
)

func Run() {
	for _, e := range conf.Config.WatchExt {
		ext.Store(strings.ToLower(e), true)
	}
	watch(true)
}

func watch(isFirst ...bool) {
	first := false
	toBuild := false
	if len(isFirst) > 0 {
		first = isFirst[0]
	}
	for _, path := range conf.Config.WatchPath {
		fs, err := utils.ScanFiles(path)
		if err != nil {
			log.Fatal(err)
			continue
		}
		for _, f := range fs {
			_, ok := ext.Load(f.Ext)
			if !ok {
				continue
			}
			modTime, exist := files.Load(f.Path)
			if !exist && !first { // 出现新增文件，则执行重建工作
				toBuild = true
				break
			}
			if !f.ModTime.Equal(modTime.(time.Time)) && !first {
				toBuild = true
				break
			}
		}
		if toBuild {
			break
		}
	}
	go rebuld()
	time.Sleep(time.Duration(conf.Config.Frequency) * time.Millisecond)
	watch()
}

func rebuld() {
	// kill current pid
	// run commands
}

func run() {

}
