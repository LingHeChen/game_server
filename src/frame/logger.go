package frame

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type CustomLogger struct {
	logger  *log.Logger
	appName string
  level int
  wd string
}

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
  colorGreen  = "\033[32m"

  DEBUG       = 1
  INFO        = 2
  WARNING     = 3
  ERROR       = 4

  DEBUG_STR   = "DEBUG"
  INFO_STR    = "INFO"
  WARNING_STR = "WARNING"
  ERROR_STR   = "ERROR"
)

var Logger *CustomLogger

func NewCustomLogger() *CustomLogger {
	return &CustomLogger{
		logger:  log.New(log.Writer(), "", 0),
	}
}

func (c *CustomLogger) InitLogger(appName string, loggerLevel int) {
  c.appName = appName
  c.level = loggerLevel
}

func ColorString(str string, color string) string {
  return color + str + colorReset
}

func getLevelString(level int) string {
  switch level {
  case DEBUG:
    return DEBUG_STR
  case INFO:
    return INFO_STR
  case WARNING:
    return WARNING_STR
  case ERROR:
    return ERROR_STR
  default:
    return "???"
  }
}

func (c *CustomLogger) Log(level int, msg string, bias int) {
  if level < c.level { return }
	fileName, lineNo, functionName := c.callerDetails(2 + bias)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
  levelString := getLevelString(level)
	color := getColor(level)
	logMessage := fmt.Sprintf(
    // 时间 | 应用名 | 等级 | 文件名(行号) | 函数名 | 信息
    "%s | %s | %s\t| %s (line %d) | %s | %s",
    currentTime,
    c.appName,
    ColorString(levelString, color),
    fileName,
    lineNo,
    functionName,
    ColorString(msg, color),
  )
	c.logger.Println(logMessage)
}

func (c *CustomLogger) callerDetails(skip int) (string, int, string) {
	pc, fileName, lineNo, ok := runtime.Caller(skip)
	if !ok {
		return "???", 0, "???"
	}

  fileName = strings.TrimPrefix(fileName, c.wd)
	fnName := runtime.FuncForPC(pc).Name()
	// Get only the function name without package path
	lastSlash := strings.LastIndex(fnName, "/")
	if lastSlash != -1 {
		fnName = fnName[lastSlash+1:]
	}
  splited := strings.Split(fnName, ".")
  fnName = strings.Join(splited[:len(splited) - 1], ".")
	return fileName, lineNo, fnName
}

func getColor(level int) string {
	switch level {
	case DEBUG:
		return colorBlue
	case INFO:
		return colorGreen // White
	case WARNING:
		return colorYellow
	case ERROR:
		return colorRed
	default:
		return colorReset
	}
}

// ===============================================================
// Log函数

func (c *CustomLogger) Debug(msg string) {
  c.Log(DEBUG, msg, 1)
}

func (c *CustomLogger) Info(msg string) {
  c.Log(INFO, msg, 1)
}

func (c *CustomLogger) Warning(msg string) {
  c.Log(WARNING, msg, 1)
}

func (c *CustomLogger) Error(msg string) {
  c.Log(ERROR, msg, 1)
}

func (c *CustomLogger) Debugf(templ string, v... any)  {
  msg := fmt.Sprintf(templ, v...)
  c.Log(DEBUG, msg, 1)
}

func (c *CustomLogger) Infof(templ string, v... any)  {
  msg := fmt.Sprintf(templ, v...)
  c.Log(INFO, msg, 1)
}

func (c *CustomLogger) Warningf(templ string, v... any)  {
  msg := fmt.Sprintf(templ, v...)
  c.Log(WARNING, msg, 1)
}

func (c *CustomLogger) Errorf(templ string, v... any)  {
  msg := fmt.Sprintf(templ, v...)
  c.Log(ERROR, msg, 1)
}

func (c *CustomLogger) Debugln(v... any)  {
  msg := fmt.Sprintln(v...)
  c.Log(DEBUG, msg, 1)
}

func (c *CustomLogger) Infoln(v... any)  {
  msg := fmt.Sprintln(v...)
  c.Log(INFO, msg, 1)
}

func (c *CustomLogger) Warningln(v... any)  {
  msg := fmt.Sprintln(v...)
  c.Log(WARNING, msg, 1)
}

func (c *CustomLogger) Errorln(v... any)  {
  msg := fmt.Sprintln(v...)
  c.Log(ERROR, msg, 1)
}


func init()  {
  Logger = NewCustomLogger()
  wd, err := os.Getwd()
  if err != nil {
    fmt.Print(err)
    os.Exit(1)
  }
  Logger.wd = wd
}
