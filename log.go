package main

import (
	"fmt"
	"log"
	"os"
)

type LogLevel int

const (
	LogFatal LogLevel = iota
	LogError
	LogWarn
	LogInfo
	LogDebug
)

func (lvl *LogLevel) UnmarshalYAML(data []byte) error {
	sData := string(data)

	if string(sData[0]) == "\"" || string(sData[0]) == "'" {
		sData = sData[1 : len(sData)-1]
	}

	switch sData {
	case "fatal":
		*lvl = LogFatal
	case "error":
		*lvl = LogError
	case "warn":
		*lvl = LogWarn
	case "info":
		*lvl = LogInfo
	case "debug":
		*lvl = LogDebug
	default:
		return fmt.Errorf("config.loglevel value invalid")
	}
	return nil
}

var logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmsgprefix)

func Log[T any](level LogLevel, msg T) {
	if level > Conf.Log {
		return
	}

	switch level {
	case LogFatal:
		logger.SetPrefix("fatal: ")
		logger.Fatalf("%v", msg)
	case LogError:
		logger.SetPrefix("error: ")
	case LogWarn:
		logger.SetPrefix("warn: ")
	case LogInfo:
		logger.SetPrefix("info: ")
	case LogDebug:
		logger.SetPrefix("debug: ")
		logger.Printf("%+v", msg)
		return
	}

	logger.Printf("%v", msg)
}
