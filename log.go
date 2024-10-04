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
		return fmt.Errorf("loglevel value invalid")
	}

	return nil
}

func (lvl *LogLevel) UnmarshalText(text []byte) error {
	return lvl.UnmarshalYAML(text)
}

var Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmsgprefix)

func Log(level LogLevel, fstring string, args ...any) {
	if level > Conf.Log {
		return
	}

	switch level {
	case LogFatal:
		Logger.SetPrefix("fatal: ")
		Logger.Fatalf(fstring, args...)
	case LogError:
		Logger.SetPrefix("error: ")
	case LogWarn:
		Logger.SetPrefix("warn: ")
	case LogInfo:
		Logger.SetPrefix("info: ")
	case LogDebug:
		Logger.SetPrefix("debug: ")
	}

	Logger.Printf(fstring, args...)
}
