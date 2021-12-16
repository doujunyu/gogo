package ff

import (
	"gogo/utility"
	"os"
)

type Log struct {
	Path string
}

// JobNewLog 初始化日志
func JobNewLog() *Log {
	root, _ := utility.UrlRootPath()
	logPath := os.Getenv("LOG_PATH")
	return &Log{
		Path: root + logPath,
	}
}

// Error 报错文件(文件名,内容)(error)
func (l *Log) Error(FileName string, content string) error {
	return l.Write("error", FileName, content)
}

// Write 写入日志 (前缀,文件名,内容)(error)
func (l *Log) Write(prefix string, FileName string, content string) error {
	return utility.LogWrite(l.Path, FileName, prefix, content) //日志内存
}
