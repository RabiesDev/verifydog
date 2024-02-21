package main

import (
	"os"
	"os/signal"
	"syscall"

	"verifydog/common"
	"verifydog/discord"
	"verifydog/server"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: false,
	})
}

func main() {
	environment := &common.Environment{}
	if err := environment.Parse(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	if environment.DevMode {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Bot is running in developer mode")
	}

	if err := common.OpenDatabase(environment.DataSourceName); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	funcBot := discord.NewBot(environment)
	if err := funcBot.Startup(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	funcServer := server.NewServer(environment, funcBot.OnAuthenticate)
	if err := funcServer.Startup(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-interruptSignal

	if err := funcBot.CloseHandle(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
