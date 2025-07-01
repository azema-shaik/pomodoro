package pomodoro

import (
	"fmt"
	"os"
	"time"

	logging "github.com/azema-shaik/logger"
)

var loggers []*logging.Logger

func getFileName() string {
	fmt.Println(os.Getwd())
	return fmt.Sprintf("logs\\%s", time.Now().Format("JAN_02_01_2006")+".log")
}

func flagType(filename string) (flag int) {
	_, err := os.Stat(filename)
	fmt.Println("err = ", err)
	if err != nil {
		return os.O_TRUNC
	}

	return os.O_APPEND

}

func NewLogger(module string) *logging.Logger {
	logger := logging.GetLogger("pomodoro." + module)
	logger.SetLevel(logging.DEBUG)
	fileName := getFileName()
	fileHandler, err := logging.GetFileHandler(fileName, os.O_CREATE|os.O_WRONLY|flagType(fileName), 0666)
	fmt.Println("err = ", err)
	fileHandler.SetLogLevel(logging.DEBUG)
	fileHandler.SetFormatter(&logging.StdFormatter{
		FormatString: "[%(asctime)s] : [%(levelname)s]: [%(lineno)d] : [%(funcName)s]: [%(msg)s]",
		DateFmt:      "Monday 02-01-2006 03:04:05 PM"})
	logger.AddHandler(fileHandler)
	loggers = append(loggers, logger)

	return logger

}

func LoggerClose() {
	for _, logger := range loggers {
		logger.Info(fmt.Sprintf("Closing logger: %s\n", logger.Name))
		logger.Close()
	}
}
