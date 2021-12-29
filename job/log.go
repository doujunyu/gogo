package job

import (
	"github.com/doujunyu/gogo/utility"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// Log 日志配置
type Log struct {
	Path    string                `Testing:"日志的路径"`
	LogChan chan *LogWriteStrings `Testing:"每次写入的日志通过管道的方式进行生成，避免并行操作"`
}

//type LogWriteFunc func(string, string, string, string)

// LogWriteStrings 将需要写入日志的信息组装
type LogWriteStrings struct {
	Url      string `Testing:"日志存放位置"`
	FileName string `Testing:"文件按名"`
	Prefix   string `Testing:"前缀"`
	Content  string `Testing:"内容"`
}

// JobNewLog 初始化日志
func NewLog() *Log {
	root, _ := utility.UrlRootPath()
	logPath := os.Getenv("LOG_PATH")
	return &Log{
		Path:    root + logPath,
		LogChan: make(chan *LogWriteStrings,1000),
	}
}

// Error 报错文件(文件名,内容)(error)
func (l *Log) Error(FileName string, content string) {
	l.Write("", FileName+"_error", content)
}

// Write 写入日志 (前缀,文件名,内容)(error)
func (l *Log) Write(prefix string, FileName string, content string) {
	l.LogChan <- &LogWriteStrings{
		Url:      l.Path,
		FileName: FileName,
		Prefix:   prefix,
		Content:  content,
	}
}

// LogWrite 写入日志 (日志路径，日志文件名，内容前缀，内容)
func LogWrite(url string, fileName string, prefix string, content string) {
	//日志内存
	path := url + "/" + fileName + "/" + time.Now().Format("2006-01-02")
	_ = os.MkdirAll(path, os.ModePerm) //生成文件
	//拼接文件绝对路径+文件名
	fileNamePath := path + "/servicing.log"
	fileInfo, err := os.Stat(fileNamePath)
	if err == nil {
		sizeG :=os.Getenv("LOG_FILE_SIZE_G")
		sizeIntG, err := strconv.ParseInt(sizeG, 10, 64)
		if err != nil {
			sizeIntG = int64(1)
		}
		if fileInfo.Size() > (sizeIntG * int64(1024576 * 1024)) {
			files, _ := ioutil.ReadDir(path)
			_ = os.Rename(fileNamePath, path+"/_"+strconv.Itoa(len(files))+".log")
		}
	}
	filePath, err := os.OpenFile(fileNamePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return
	}
	defer func(filePath *os.File) {
		err := filePath.Close()
		if err != nil {

		}
	}(filePath)
	//写入文件
	logger := log.New(filePath, "["+prefix+"]:", log.Lmicroseconds)
	logger.Printf(content)
	return
}
