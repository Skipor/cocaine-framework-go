package cocaine

import (
	"errors"
	"fmt"
	"time"
)

type Logger struct {
	socketWriter
	verbosity int
}

const (
	LOGINGNORE = iota
	LOGERROR
	LOGWARN
	LOGINFO
	LOGDEBUG
)

func NewLogger() (logger *Logger, err error) {
	temp, err := NewService("logging", "localhost", 10053)
	if err != nil {
		return
	}
	defer temp.Close()

	res := <-temp.Call("verbosity")
	if res.Err() != nil {
		err = errors.New("Unable to receive verbosity")
		return
	}
	var verbosity int = 0
	if err = res.Extract(&verbosity); err != nil {
		return
	}

	sock, err := NewWSocket("tcp", temp.ResolveResult.AsString(), time.Second*5)
	if err != nil {
		return
	}

	//Create logger
	logger = &Logger{sock, verbosity}
	return
}

// Blocked
func (logger *Logger) log(level int64, message ...interface{}) bool {
	msg := ServiceMethod{MessageInfo{0, 0}, []interface{}{level, fmt.Sprintf("app/%s", flag_app), fmt.Sprint(message...)}}
	logger.Write() <- Pack(&msg)
	return true
}

func (logger *Logger) Err(message ...interface{}) {
	_ = LOGERROR <= logger.verbosity && logger.log(LOGERROR, message...)
}

func (logger *Logger) Warn(message ...interface{}) {
	_ = LOGWARN <= logger.verbosity && logger.log(LOGWARN, message...)
}

func (logger *Logger) Info(message ...interface{}) {
	_ = LOGINFO <= logger.verbosity && logger.log(LOGINFO, message...)
}

func (logger *Logger) Debug(message ...interface{}) {
	_ = LOGDEBUG <= logger.verbosity && logger.log(LOGDEBUG, message...)
}

func (logger *Logger) Close() {
	logger.socketWriter.Close()
}
