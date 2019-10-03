package serv

import (
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
	files      = &sync.Map{}
	filesCount = 0
	ext        = &sync.Map{}
	cmd        *exec.Cmd
)

func Run() {
	for _, e := range conf.Config.WatchExt {
		ext.Store(strings.ToLower(e), true)
	}
	watch()
}

// 文件少了(删除了)或者多了(新增了)或者发生更改，都要触发
func watch() {
	tmpCount := 0
	sholdRestart := false
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
			tmpCount++
			modTime, ok := files.Load(f.Path)
			if ok && !f.ModTime.Equal(modTime.(time.Time)) {
				sholdRestart = true
				continue
			}
			if !ok {
				sholdRestart = true
				files.Store(f.Path, f.ModTime)
			}
		}
	}
	if tmpCount != filesCount {
		filesCount = tmpCount
		sholdRestart = true
	}
	if sholdRestart {
		go restart()
	}
	time.Sleep(time.Duration(conf.Config.Frequency) * time.Millisecond)
	watch()
}

func restart() {
	kill()
	for _, command := range conf.Config.Commands {
		if len(command) > 0 {
			if err := execCommand(command[0], command[1:]...); err != nil {
				if conf.Config.Strict {
					break
				}
			}
		}
	}
}

func execCommand(commandName string, params ...string) error {
	cmd = exec.Command(commandName, params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Kill kills the running command process
func kill() {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
		}
	}()
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			log.Println(err)
		}
	}
}
