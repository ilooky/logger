package main

import (
	"github.com/ilooky/logger"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger(logger.Config{Level: "debug", Release: false, Path: "", Style: "json"})
	logger.DebugKV("not find", "consul", "192.168.1.2")
	logger.Debug("this is log out")
	logger.Info("this is info")
	logger.Warnf("this is warn:%s", "warning")
	logger.Error("this is error")
	logger.ErrorKv("this is zap error")
	logger.ErrorKv("this is zap error", zap.String("user", "us"))
}
