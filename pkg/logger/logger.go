package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	info *log.Logger
	warn *log.Logger
	err  *log.Logger
}

var logger *Logger

func SetLogger(l *Logger) {
	logger = l
}

func NewFileAndStdOutLogger(fn string) (*Logger, error) {
	file, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return New(os.Stdout, file)
}

func New(writers ...io.Writer) (*Logger, error) {
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	fn := os.Getenv("LOG_FILE")
	if fn == "" {
		fn = "app.log"
	}

	mw := io.MultiWriter(writers...)
	return &Logger{
		info: log.New(mw, "INFO: ", log.LstdFlags),
		warn: log.New(mw, "WARN: ", log.LstdFlags),
		err:  log.New(mw, "ERROR: ", log.LstdFlags),
	}, nil
}

func Info(v ...interface{}) {
	logger.info.Println(v...)
}

func Warn(v ...interface{}) {
	logger.warn.Println(v...)
}

func Error(v ...interface{}) {
	logger.err.Println(v...)
}
