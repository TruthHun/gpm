package utils

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

//返回的目录扫描结果
type File struct {
	IsDir   bool      //是否是目录
	Path    string    //文件路径
	Ext     string    //文件扩展名
	Name    string    //文件名
	Size    int64     //文件大小
	ModTime time.Time //文件修改时间戳
}

//目录扫描
//@param			dir			需要扫描的目录
//@return			fl			文件列表
//@return			err			错误
func ScanFiles(dir string) (files []File, err error) {
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			path = strings.Replace(path, "\\", "/", -1) //文件路径处理
			files = append(files, File{
				IsDir:   info.IsDir(),
				Path:    path,
				Ext:     strings.ToLower(filepath.Ext(path)),
				Name:    info.Name(),
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
		}
		return err
	})
	return
}
