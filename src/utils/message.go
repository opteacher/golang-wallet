package utils

import (
	"fmt"
	"log"
	"strconv"
	"reflect"
	"errors"
	"os"
	"strings"
	"time"
	"path"
	"runtime"
)

const (
	ERROR = iota
	WARNING
	INFO
	DEBUG
)

var levels = map[int]string {
	ERROR: "E", WARNING: "W", INFO: "I", DEBUG: "D",
}

var logFiles = make(map[string]*os.File)

func GetIdxMsg(idx string) string {
	msgSet := GetConfig().GetMsgsSettings()
	var msgs = map[string]string {}
	switch idx[0:1] {
	case "E":
		msgs = msgSet.Errors
	case "W":
		msgs = msgSet.Warnings
	case "I":
		msgs = msgSet.Information
	case "D":
		msgs = msgSet.Debugs
	}
	return msgs[idx]
}

func LogIdxEx(level int, id int, detail ...interface {}) error {
	msgSet := GetConfig().GetMsgsSettings()
	if !msgSet.Logs.Debug && level == DEBUG {
		return procsDetail(detail)
	}

	log.SetFlags(log.Lshortfile)

	var msgs = map[string]string {}
	var prefix string
	switch level {
	case ERROR:
		msgs = msgSet.Errors
		prefix = "E"
	case WARNING:
		msgs = msgSet.Warnings
		prefix = "W"
	case INFO:
		msgs = msgSet.Information
		prefix = "I"
	case DEBUG:
		msgs = msgSet.Debugs
		prefix = "D"
	}

	if msg, ok := msgs[prefix + fmt.Sprintf("%04d", id)]; ok {
		return logMsgEx(level, msg, detail...)
	} else {
		return logMsgEx(level, "未找到指定的消息信息：%d", id)
	}
}

func logMsgEx(level int, msg string, detail ...interface {}) error {
	msgSet := GetConfig().GetMsgsSettings()
	if !msgSet.Logs.Debug && level == DEBUG {
		return procsDetail(detail)
	}

	log.SetFlags(log.Lshortfile)

	if _, ok := levels[level]; !ok {
		level = ERROR
	}

	if len(detail) >= 1 && detail[0] != nil {
		msg = fmt.Sprintf(msg, detail...)
	}

	strLevel := GetConfig().GetMsgsSettings().Level[strconv.Itoa(level)]
	logTxt := fmt.Sprintf("[%s] %s\n", strLevel, msg)

	var err error
	fileLog := GetConfig().msgs.Storage.File
	if fileLog.Active {
		// 如果文件夹不存在，创建
		if len(fileLog.Path) != 0 {
			if _, err = os.Stat(fileLog.Path); os.IsNotExist(err) {
				if err = os.Mkdir(fileLog.Path, os.ModePerm); err != nil {
					log.Fatal(err)
				}
			}
		}
		// 添加发生的文件和行数
		location := "???:?? "
		if _, file, line, ok := runtime.Caller(2); ok {
			projPath := ""
			if projPath, err = os.Getwd(); err != nil {
				log.Fatal(err)
			}
			projPath = path.Join(projPath, "src")
			file = strings.Replace(file, projPath, "", -1)
			location = fmt.Sprintf("%s:%d ", file, line)
		}
		// 根据split做分离
		switch fileLog.Split {
		case "type":
		case "level":
			fallthrough
		default:
			// 组装文件名
			logFileTag := levels[level]
			if _, ok := logFiles[logFileTag]; !ok {
				splitName := GetConfig().msgs.Level[strconv.Itoa(level)]
				fileName := genLogFileName(fileLog.NameFormat, splitName)
				if logFiles[logFileTag], err = os.OpenFile(
					path.Join(fileLog.Path, fileName),
					os.O_APPEND | os.O_CREATE | os.O_WRONLY,
					0644,
				); err != nil {
					log.Fatal(err)
				}
			}
			if _, err = logFiles[logFileTag].Write([]byte(location + logTxt)); err != nil {
				log.Fatal(err)
			}
		}
	}

	log.Output(3, logTxt)
	if len(detail) > 1 {
		return procsDetail(msg)
	} else {
		return procsDetail(detail[0])
	}
}

func genLogFileName(nameFmt string, split string) string {
	ret := strings.Replace(nameFmt, "{split}", split, -1)
	ret = strings.Replace(ret, "{time}", time.Now().Format("2006-01-02"), -1)
	ret = strings.Replace(ret, "{suffix}", "", -1)
	return ret
}

func procsDetail(detail interface {}) error {
	if detail != nil {
		typName := reflect.TypeOf(detail).Name()
		switch typName {
		case "error":
			return detail.(error)
		case "string":
			return errors.New(detail.(string))
		case "int":
			fallthrough
		case "int32":
			fallthrough
		case "int64":
			fallthrough
		case "uint":
			fallthrough
		case "uint32":
			fallthrough
		case "uint64":
			return errors.New(fmt.Sprintf("%d", detail))
		default:
			return nil
		}
	}
	return nil
}

func LogMsgEx(level int, msg string, detail ...interface {}) error {
	return logMsgEx(level, msg, detail...)
}

func CloseAllLogStorage() {
	for _, f := range logFiles {
		f.Close()
	}
}