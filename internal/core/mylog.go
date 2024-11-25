package core

import (
	"os"
	"sync"
	"time"
)

type MyLogger struct {
	dscrpt *os.File
	mutex  sync.Mutex
}

func NewLogger() *MyLogger {
	return &MyLogger{dscrpt: nil, mutex: sync.Mutex{}}
}

func (logger *MyLogger) Close() {
	logger.mutex.Lock()
	err := logger.dscrpt.Close()
	if err != nil {
		os.Stderr.WriteString("error log close: " + err.Error())
		logger.LogERROR(err.Error())
	}
	logger.mutex.Unlock()

}

func (logger *MyLogger) LogINFO(msg_info string) {
	logger.mutex.Lock()
	_, err := logger.dscrpt.WriteString("INFO, " + time.Now().Format("2006-01-02 15:04:05.000") + ", " + msg_info + "\n")
	if err != nil {
		os.Stderr.WriteString("error logging: " + err.Error())
	}
	logger.mutex.Unlock()
}
func (logger *MyLogger) LogWARN(msg_info string) {
	logger.mutex.Lock()
	_, err := logger.dscrpt.WriteString("WARN, " + time.Now().Format("2006-01-02 15:04:05.000") + ", " + msg_info + "\n")
	if err != nil {
		os.Stderr.WriteString("error logging: " + err.Error())
	}
	logger.mutex.Unlock()
}
func (logger *MyLogger) LogERROR(msg_info string) {
	logger.mutex.Lock()
	_, err := logger.dscrpt.WriteString("ERROR, " + time.Now().Format("2006-01-02 15:04:05.000") + ", " + msg_info + "\n")
	if err != nil {
		os.Stderr.WriteString("error logging: " + err.Error())
	}
	logger.mutex.Unlock()
}
func (logger *MyLogger) LogFATAL(msg_info string) {
	logger.mutex.Lock()
	_, err := logger.dscrpt.WriteString("FATAL, " + time.Now().Format("2006-01-02 15:04:05.000") + ", " + msg_info + "\n")
	if err != nil {
		os.Stderr.WriteString("error logging: " + err.Error())
	}
	os.Exit(1)
	logger.mutex.Unlock()
}

func (logger *MyLogger) SetDescriptor(dscrpt *os.File) {
	logger.dscrpt = dscrpt
}
