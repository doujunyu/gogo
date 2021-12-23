package ff

import (
	utility "gogo/utility"
	"os"
)

type Log struct {
	Path string
	LogChan *chan utility.LogWriteStrings
}

// JobNewLog 初始化日志
func JobNewLog(logChan *chan utility.LogWriteStrings) *Log {
	root, _ := utility.UrlRootPath()
	logPath := os.Getenv("LOG_PATH")
	return &Log{
		Path: root + logPath,
		LogChan: logChan,
	}
}

// Error 报错文件(文件名,内容)(error)
func (l *Log) Error(FileName string, content string) {
	l.Write("", FileName+"_error", content)
}

// Write 写入日志 (前缀,文件名,内容)(error)
func (l *Log) Write(prefix string, FileName string, content string) {
	 *l.LogChan <- utility.LogWriteStrings{
		 Url: l.Path,
		 FileName: FileName,
		 Prefix: prefix,
		 Content: content,
	 }

}
