package utils

import (
	"fmt"
	"log"
	"strconv"
	"reflect"
	"errors"
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

func LogIdxEx(level int, id int, detail interface {}) error {
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
		return logMsgEx(level, msg, detail)
	} else {
		return logMsgEx(level, "未找到指定的消息信息：%d", id)
	}
}

func logMsgEx(level int, msg string, detail interface {}) error {
	msgSet := GetConfig().GetMsgsSettings()
	if !msgSet.Logs.Debug && level == DEBUG {
		return procsDetail(detail)
	}

	log.SetFlags(log.Lshortfile)

	if _, ok := levels[level]; !ok {
		level = ERROR
	}

	if detail != nil {
		msg = fmt.Sprintf(msg, detail)
	}

	strLevel := GetConfig().GetMsgsSettings().Level[strconv.Itoa(level)]
	log.Output(3, fmt.Sprintf("[%s] %s\n", strLevel, msg))
	return procsDetail(detail)
}

func procsDetail(detail interface {}) error {
	if detail != nil {
		switch reflect.TypeOf(detail).Name() {
		case "error":
			return detail.(error)
		case "string":
			return errors.New(detail.(string))
		default:
			return nil
		}
	}
	return nil
}

func LogMsgEx(level int, msg string, detail interface {}) error {
	return logMsgEx(level, msg, detail)
}