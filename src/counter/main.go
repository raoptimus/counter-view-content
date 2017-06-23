package main

import (
	"github.com/raoptimus/gserv/config"
	"github.com/raoptimus/gserv/service"
	"github.com/raoptimus/rlog"
)

var log *rlog.Logger
var counter *Counter

func main() {
	if service.Exists() {
		panic("Service is already running. The App will be close")
	}

	flags := rlog.LOG_ERR | rlog.LOG_CRIT | rlog.LOG_EMERG
	//flags |= rlog.LOG_INFO | rlog.LOG_WARNING | rlog.LOG_DEBUG | rlog.LOG_NOTICE

	var err error
	log, err = rlog.NewLogger(rlog.LoggerTypeStd, "", flags)
	if err != nil {
		panic(err)
	}

	service.Init(&service.BaseService{
		Logger:   log,
		Start:    start,
		Stop:     stop,
		Location: service.GetTimeMoskow(),
	})

	service.StartProfiler(":8061")
	service.Start(true)
}

func start() {
	Init()
	counter = NewCounter()
	addr := config.String("ServerAddr", "/tmp/content-counter.sock")
	NewServer(addr)
}

func stop() {
	//todo backup counter mem
}
