package utils

import (
	"fmt"
	"log"
	"strconv"
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

var enableDebugLog = false
func EnableDebugLog(enable bool) {
	enableDebugLog = enable
}

func LogIdxEx(level int, id int, detail interface {}) error {
	if !enableDebugLog && level == DEBUG {
		return detail.(error)
	}

	log.SetFlags(log.Lshortfile)

	msgSet := GetConfig().GetMsgsSettings()
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

	if msg, ok := msgs[prefix + strconv.Itoa(id)]; ok {
		return LogMsgEx(level, msg, detail)
	} else {
		return LogMsgEx(level, "未找到指定的消息信息：%d", id)
	}
}

func LogMsgEx(level int, msg string, detail interface {}) error {
	if !enableDebugLog && level == DEBUG {
		return detail.(error)
	}

	log.SetFlags(log.Lshortfile)

	if _, ok := levels[level]; !ok {
		level = ERROR
	}

	if detail != nil {
		msg = fmt.Sprintf(msg, detail)
	}

	strLevel := GetConfig().GetMsgsSettings().Level[strconv.Itoa(level)]
	log.Output(0, fmt.Sprintf("[%s] %s\n", strLevel, msg))
	if level == ERROR {
		detail = errors.New(msg)
	}
	return detail.(error)
}