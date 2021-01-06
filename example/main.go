package main

import (
	"github.com/ilooky/logger"
	"time"
)

func main() {
	logger.InitLogger("diagram.log", "debug")
	logger.DebugKV("not find", "consul", "192.168.1.2")
	logger.Debug("this is log out")
	time.Sleep(3 * time.Second)
}
