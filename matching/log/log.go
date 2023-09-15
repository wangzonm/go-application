package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LEVEL byte

const (
	DEBUG LEVEL = iota
	INFO
	WARN
	ERROR
)

type FileLogger struct {
	fileDir        string        // 日志文件保存的目录
	fileName       string        // 日志文件名（无需包含日期和扩展名）
	prefix         string        // 日志消息的前缀
	logLevel       LEVEL         // 日志等级
	logFile        *os.File      // 日志文件
	date           *time.Time    // 日志当前日期
	lg             *log.Logger   // 系统日志对象
	mu             *sync.RWMutex // 读写锁，在进行日志分割和日志写入时需要锁住
	logChan        chan string   // 日志消息通道，以实现异步写日志
	stopTickerChan chan bool     // 停止定时器的通道
}

const DATE_FORMAT = "2006-01-02"

var fileLogger *FileLogger

func Init(fileDir, fileName, prefix string, level LEVEL) error {
	CloseLogger()

	f := &FileLogger{
		fileDir:        fileDir,
		fileName:       fileName,
		prefix:         prefix,
		mu:             new(sync.RWMutex),
		logChan:        make(chan string, 5000),
		stopTickerChan: make(chan bool, 1),
	}

	switch strings.ToUpper(string(level)) {
	case "DEBUG":
		f.logLevel = DEBUG
	case "WARN":
		f.logLevel = WARN
	case "ERROR":
		f.logLevel = ERROR
	default:
		f.logLevel = INFO
	}

	t, _ := time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))
	f.date = &t

	//f.isExistOrCreateFileDir()

	fullFileName := filepath.Join(f.fileDir, f.fileName+".log")
	file, err := os.OpenFile(fullFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	f.logFile = file

	f.lg = log.New(f.logFile, prefix, log.LstdFlags|log.Lmicroseconds)

	go f.logWriter()
	go f.fileMonitor()

	fileLogger = f

	return nil
}

func CloseLogger() {
	if fileLogger != nil {
		fileLogger.stopTickerChan <- true
		close(fileLogger.stopTickerChan)
		close(fileLogger.logChan)
		fileLogger.lg = nil
		fileLogger.logFile.Close()
	}
}

func (f *FileLogger) logWriter() {
	defer func() { recover() }()

	for {
		str, ok := <-f.logChan
		if !ok {
			return
		}

		f.mu.RLock()
		f.lg.Output(2, str)
		f.mu.RUnlock()
	}
}

func (f *FileLogger) fileMonitor() {
	defer func() { recover() }()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if f.isMustSplit() {
				if err := f.split(); err != nil {
					Error("Log split error: %v\n", err)
				}
			}
		case <-f.stopTickerChan:
			return
		}
	}
}

func (f *FileLogger) isMustSplit() bool {
	t, _ := time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))
	return t.After(*f.date)
}

func (f *FileLogger) split() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	logFile := filepath.Join(f.fileDir, f.fileName)
	logFileBak := logFile + "-" + f.date.Format(DATE_FORMAT) + ".log"

	if f.logFile != nil {
		f.logFile.Close()
	}

	err := os.Rename(logFile, logFileBak)
	if err != nil {
		return err
	}

	t, _ := time.Parse(DATE_FORMAT, time.Now().Format(DATE_FORMAT))
	f.date = &t

	f.logFile, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	f.lg = log.New(f.logFile, f.prefix, log.LstdFlags|log.Lmicroseconds)

	return nil
}

func Debug(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLogger.logLevel <= DEBUG {
		fileLogger.logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[DEBUG]"+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLogger.logLevel <= INFO {
		fileLogger.logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[INFO]"+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	if fileLogger.logLevel <= ERROR {
		fileLogger.logChan <- fmt.Sprintf("[%v:%v]", filepath.Base(file), line) + fmt.Sprintf("[ERROR]"+format, v...)
	}
}
