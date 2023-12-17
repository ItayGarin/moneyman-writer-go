package logger

import "go.uber.org/zap"

var x *zap.SugaredLogger
var isInit bool = false

func Init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	x = logger.Sugar()
	isInit = true
}

func InitDev() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	x = logger.Sugar()
	isInit = true
}

func Logger() *zap.SugaredLogger {
	if !isInit {
		Init()
		isInit = true
	}

	return x
}
