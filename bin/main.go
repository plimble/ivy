package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
	"github.com/plimble/ivy"
)

func main() {
	config, err := ivy.GetConfig()
	if err != nil {
		panic(err)
	}

	server, err := ivy.NewServer(config)
	if err != nil {
		panic(err)
	}

	log.Infof("Iris %s Running at %s", iris.Version, config.Addr)

	server.Run(iris.Addr(config.Addr), iris.WithoutStartupLog)
}
