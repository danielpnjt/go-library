package log

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	httpRequest      = "REQUEST"
	httpResponse     = "RESPONSE"
	folder           = "logs"
	timeformat       = "2006-01-02T15:04:05-0700"
	nameformat       = "log-2006-01-02.log"
	nameformatTrxLog = "trxlog-2006-01-02.log"
)

var (
	currentFile     *os.File
	logText         *logrus.Logger
	logJSON         *logrus.Logger
	currentFileName string
	err             error
)

func Init(serviceName string, debug bool) {
	SetText()
	SetJSON()
	SetFolder()

	if err != nil {
		fmt.Println(err)
	}

	if debug {
		logText.SetLevel(logrus.DebugLevel)
		logJSON.SetLevel(logrus.DebugLevel)
	} else {
		logText.SetLevel(logrus.InfoLevel)
		logJSON.SetLevel(logrus.InfoLevel)
	}
}

func SetFolder() {
	dir, _ := os.Getwd()
	folderlogs := dir + "/" + folder

	if _, err := os.Stat(folderlogs); os.IsNotExist(err) {
		err := os.Mkdir(folderlogs, 0777)
		fmt.Println(err)
	}
}

func SetJSON() {
	logJSON = logrus.New()
	formatter := new(logrus.JSONFormatter)
	formatter.DisableTimestamp = true
	logJSON.SetFormatter(formatter)
}

func SetText() {
	logText = logrus.New()
	formatter := new(logrus.TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableQuote = true
	logText.SetFormatter(formatter)
}

func SetLogFile(mode int) string {
	currentTime := time.Now()
	timestamp := currentTime.Format(timeformat)

	fileFormat := nameformat

	if mode == 1 {
		fileFormat = nameformatTrxLog
	}

	filename := folder + "/" + currentTime.Format(fileFormat)
	if filename == currentFileName {
		return timestamp
	}

	newLogFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	} else {
		currentFileName = filename
		logText.SetOutput(newLogFile)
		logJSON.SetOutput(newLogFile)

		if currentFile != nil {
			currentFile.Close()
			currentFile = newLogFile
		}
	}

	return timestamp
}

func LogDebug(msg string) {
	timestamp := SetLogFile(0)
	logText.Debug(fmt.Sprintf("%s [%s] %s", timestamp, "", msg))
}

func Minify(r interface{}) map[string]interface{} {
	js, _ := json.Marshal(r)
	var m map[string]interface{}
	_ = json.Unmarshal(js, &m)

	minifyThreshold := 100

	for k, v := range m {
		if k == "response_data" || k == "responseData" {
			_, ok := v.(map[string]interface{})
			if !ok {
				m[k] = map[string]interface{}{}
			}
		}

		s := fmt.Sprintf("%v", v)
		if len(s) > minifyThreshold {
			m[k] = "panjang"

			_, ok := v.(string)
			if !ok || k == "response_data" || k == "responseData" {
				m[k] = map[string]interface{}{}
			}
		}
	}

	jsm, _ := json.Marshal(m)
	strJsm := string(jsm)

	_ = json.Unmarshal([]byte(strJsm), &m)

	return m
}
