package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/unification-com/unode-onboard-api/pkg/handlers"
	db "github.com/unification-com/unode-onboard-api/pkg/models"
)

func main() {
	if err := godotenv.Load("dev.env"); err != nil {
		panic(err.Error())
	}

	logrus.SetLevel(logrus.DebugLevel)

	dbClient := db.NewDBClient()

	defer dbClient.CloseConnection()

	server := handlers.NewServer(dbClient)

	// syscall.SIGUSR1, syscall.SIGUSR2 will only work in linux
	// open the unode api in linux os or wsl to have access to SIGUSR1 and SIGUSR2 signals
	// Alternative: Comment it out and run it go with restart approach to change log level
	// Commands to run it:
	// kubectl exec <my-app-pod> -- kill -s SIGUSR1 <PID>
	// kubectl exec <my-app-pod> -- kill -s SIGUSR2 <PID>
	// example to get PID of the process: kubectl exec -it unode-api-8cb75cdc5-bvzsc -- ps aux (Lets say PID is 1)
	// example to change log level: kubectl exec unode-api-8cb75cdc5-bvzsc -- kill -s SIGUSR1 1 (Changes loglevel to debug)

	// TODO: build with linux os or wsl
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGUSR1, syscall.SIGUSR2)
	// go func() {
	// 	for {
	// 		sig := <-sigChan
	// 		switch sig {
	// 		case syscall.SIGUSR1:
	// 			logrus.SetLevel(logrus.DebugLevel)
	// 			logrus.Warnln("Log level changed to DEBUG")
	// 		case syscall.SIGUSR2:
	// 			logrus.SetLevel(logrus.InfoLevel)
	// 			logrus.Warnln("Log level changed to INFO")
	// 		}
	// 	}
	// }()

	go func() {
		if err := server.RunServer(); err != nil {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	logrus.Info("Closing the Server")
}
