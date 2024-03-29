package gogo_log

import (
	_ "github.com/joho/godotenv/autoload"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var GlobalLogData *Log
var GlobalLogOnce sync.Once

// Log 日志配置
type Log struct {
	Path    string                `Testing:"日志的路径"`
	LogChan chan *LogWriteStrings `Testing:"每次写入的日志通过管道的方式进行生成，避免并行操作"`
}

// LogWriteStrings 将需要写入日志的信息组装
type LogWriteStrings struct {
	Url      string    `Testing:"日志存放位置"`
	FileName string    `Testing:"文件按名"`
	Prefix   string    `Testing:"前缀"`
	Content  string    `Testing:"内容"`
	Ctime    time.Time `Testing:"时间"`
}

func init() {
	GlobalLogOnce.Do(func() {
		dir, _ := os.Getwd()
		logPath := os.Getenv("LOG_PATH")
		GlobalLogData = &Log{
			Path:    strings.Replace(dir, "\\", "/", -1) + logPath,
			LogChan: make(chan *LogWriteStrings, 1000),
		}
	})
}

// Error 报错文件(文件名,内容)(error)
func Error(FileName string, prefix string, content string) {
	Write("error_"+FileName, prefix, content)
}

// Write 写入日志 (前缀,文件名,内容)(error)
func Write(FileName string, prefix string, content string) {
	GlobalLogData.LogChan <- &LogWriteStrings{
		Url:      GlobalLogData.Path,
		FileName: FileName,
		Prefix:   prefix,
		Content:  content,
		Ctime:    time.Now(),
	}
}

func LogChanOut() {
	for {
		data := <-GlobalLogData.LogChan
		if data != nil {
			logWrite(data.Url, data.FileName, data.Prefix, data.Content, data.Ctime) //日志内存
		}
	}
}

// logWrite 写入日志 (日志路径，日志文件名，内容前缀，内容)
func logWrite(url string, fileName string, prefix string, content string, ctime time.Time) {
	//日志内存
	path := url + "/" + fileName + "/" + ctime.Format("2006-01-02")
	_ = os.MkdirAll(path, os.ModePerm) //生成文件
	//拼接文件绝对路径+文件名
	fileNamePath := path + "/servicing.gogo_log"
	fileInfo, err := os.Stat(fileNamePath)
	if err == nil {
		sizeG := os.Getenv("LOG_FILE_SIZE_G")
		sizeIntG, err := strconv.ParseInt(sizeG, 10, 64)
		if err != nil {
			sizeIntG = int64(1)
		}
		if fileInfo.Size() > (sizeIntG * int64(1024576*1024)) {
			files, _ := ioutil.ReadDir(path)
			_ = os.Rename(fileNamePath, path+"/_"+strconv.Itoa(len(files))+".gogo_log")
		}
	}
	filePath, err := os.OpenFile(fileNamePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return
	}
	defer filePath.Close()
	//写入文件
	logger := log.New(filePath, "["+prefix+"]:", log.Lmicroseconds)
	logger.Printf("|" + ctime.Format("15:04:05") + "|" + content)
	return
}
