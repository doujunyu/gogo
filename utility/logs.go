package utility

import (
	"log"
	"os"
	"time"
)

// Write 写入日志 (日志路径，日志文件名，内容前缀，内容)
func LogWrite(url string, FileName string, prefix string, content string) error {
	//日志内存
	path := url + FileName
	_ = os.MkdirAll(path, os.ModePerm) //判断文件夹是否存在
	//拼接文件绝对路径+文件名
	fileNamePath := path + "/" + time.Now().Format("2006-01-02-15") + ".log"
	filePath, err := os.OpenFile(fileNamePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer filePath.Close()
	//写入文件
	logger := log.New(filePath, "["+prefix+"]:", log.Llongfile)
	logger.Printf(content)
	return nil
}
