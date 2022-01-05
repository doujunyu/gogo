package log

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
	Url      string `Testing:"日志存放位置"`
	FileName string `Testing:"文件按名"`
	Prefix   string `Testing:"前缀"`
	Content  string `Testing:"内容"`
}

func init()  {
	GlobalLogOnce.Do(func() {
		dir, _ := os.Getwd()
		logPath := os.Getenv("LOG_PATH")
		GlobalLogData =  &Log{
			Path:    strings.Replace(dir, "\\", "/", -1) + logPath,
			LogChan: make(chan *LogWriteStrings,1000),
		}
	})
}


// Error 报错文件(文件名,内容)(error)
func Error(FileName string, content string) {
	Write("", FileName+"_error", content)
}
// Write 写入日志 (前缀,文件名,内容)(error)
func Write(prefix string, FileName string, content string) {
	GlobalLogData.LogChan <- &LogWriteStrings{
		Url:      GlobalLogData.Path,
		FileName: FileName,
		Prefix:   prefix,
		Content:  content,
	}
}

func LogChanOut() {
	for {
		data := <-GlobalLogData.LogChan
		logWrite(data.Url, data.FileName, data.Prefix, data.Content) //日志内存
	}
}

// logWrite 写入日志 (日志路径，日志文件名，内容前缀，内容)
func logWrite(url string, fileName string, prefix string, content string) {
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
