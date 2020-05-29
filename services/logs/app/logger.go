package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Logger interface {
	Access(timestamp int64, serverName, method, requestId, user, duration string)
	Error(timestamp int64, serverName string, requestId string, data string)
	System(timestamp int64, serverName string, data string)
}

type LoggerData struct {
	oldName string
	writer  io.WriteCloser
	Path    string
}

func InitLogger(path string) Logger {
	return &LoggerData{
		Path: path,
	}
}

func (l *LoggerData) getFileName() string {
	y, m, d := time.Now().Date()
	return filepath.Join(l.Path, fmt.Sprintf("%v.%v.%v.log", d, m, y))
}

func (l *LoggerData) setNewWriter() {

	if l.writer != nil {
		_ = l.writer.Close()
	}
	l.oldName = l.getFileName()
	l.writer, _ = os.OpenFile(l.oldName, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
}

var accessTemplate = "%8s %22s |  %12s  |  %28s  |  %12s  |  %v (%v)\n"

func (l *LoggerData) Access(timestamp int64, serverName, method, requestId, user, duration string) {

	if l.getFileName() != l.oldName {
		l.setNewWriter()
	}

	if os.Getenv("APP_MODE") == "debug" {
		fmt.Print(fmt.Sprintf(accessTemplate,
			"(ACCESS)",
			timestampToDate(timestamp),
			requestId,
			user,
			duration,
			serverName,
			method,))
	}

	_, _ = l.writer.Write([]byte(fmt.Sprintf(accessTemplate,
		"(ACCESS)",
		timestampToDate(timestamp),
		requestId,
		user,
		duration,
		serverName,
		method, )))
}

var errorTemplate = "%8s %22s |  %12s  |  %28s  |  %v\n"

func (l *LoggerData) Error(timestamp int64, serverName string, requestId string, data string) {
	if l.getFileName() != l.oldName {
		l.setNewWriter()
	}

	if os.Getenv("APP_MODE") == "debug" {
		fmt.Print(fmt.Sprintf(errorTemplate, "(ERROR)", timestampToDate(timestamp), requestId, serverName, data))
	}

	_, _ = l.writer.Write([]byte(
		fmt.Sprintf(
			errorTemplate, "(ERROR)", timestampToDate(timestamp), requestId, serverName, data,
		)))
}

var systemTemplate = "%8s %22s |  %12s  |  %v\n"

func (l *LoggerData) System(timestamp int64, serverName string, data string) {

	if l.getFileName() != l.oldName {
		l.setNewWriter()
	}

	if os.Getenv("APP_MODE") == "debug" {
		fmt.Print(fmt.Sprintf(systemTemplate, "(SYSTEM)", timestampToDate(timestamp), serverName, data))
	}

	_, _ = l.writer.Write([]byte(fmt.Sprintf(systemTemplate, "(SYSTEM)", timestampToDate(timestamp), serverName, data)))
}

func timestampToDate(in int64) string {
	return time.Unix(0, in).In(time.Local).Format(time.RFC822)
}
