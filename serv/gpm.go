package serv

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/TruthHun/gpm/conf"
	"github.com/TruthHun/gpm/utils"
)

var (
	files     sync.Map
	ext       sync.Map
	processes []*os.Process
)

func Run() {
	for _, e := range conf.Config.WatchExt {
		ext.Store(strings.ToLower(e), true)
	}
	watch(true)
}

// 文件少了(删除了)或者多了(新增了)或者发生更改，都要触发
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

			if !first && !exist { // 出现新增文件，则执行重建工作
				toBuild = true
			}
			if !first && modTime != nil && !f.ModTime.Equal(modTime.(time.Time)) {
				toBuild = true
			}
			if toBuild {
				files.Store(f.Path, f.ModTime)
			}
		}
		if toBuild {
			break
		}
	}
	if toBuild {
		go rebuld()
	}
	time.Sleep(time.Duration(conf.Config.Frequency) * time.Millisecond)
	watch()
}

func rebuld() {
	for _, process := range processes {
		if process != nil {
			process.Kill()
		}
	}
	if len(conf.Config.Commands) > 0 {
		execCommand(conf.Config.Commands[0], conf.Config.Commands[1:]...)
	}
}

func execCommand(commandName string, params ...string) bool {
	cmd := exec.Command(commandName, params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		return false
	}

	cmd.Start()

	processes = append(processes, cmd.Process)

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Printf("[%v] %v", time.Now().Local().Format("15:04:05"), line)
	}

	cmd.Wait()
	return true
}
