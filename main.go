package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"caroline-weisberg.fun/iljournalierserver/errors"
	"caroline-weisberg.fun/iljournalierserver/intake"
	"caroline-weisberg.fun/iljournalierserver/utils"
)

func main() {
	utils.AlwaysLog("IlJournalierServer: start!")

	intakeConfig, err := intake.ReadIntakeConfiguration()
	if err != nil {
		utils.AlwaysLog(errors.J(err, "Read Intake Configuration failed").Error())
		return
	}

	ctx := context.Background()

	diContainer, err := newDIContainer(ctx, intakeConfig)

	if err != nil {
		utils.AlwaysLog(errors.J(err, "create di container failed").Error())
		return
	}

	mainRouter := newMainRouter(diContainer)

	mux := http.NewServeMux()
	mux.Handle("/", &mainRouter)

	server := http.Server{Handler: mux}

	listener, err := net.Listen(intakeConfig.Transport.Network, intakeConfig.Transport.Address)
	if err != nil {
		utils.AlwaysLog(errors.J(err, "create listener failed").Error())
		return
	}

	go func() {
		utils.AlwaysLog("IlJournalierServer: will listen", intakeConfig.Transport.Network, intakeConfig.Transport.Address)
		err := server.Serve(listener)
		if err != nil {
			utils.AlwaysLog("IlJournalierServer: returned", err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		utils.AlwaysLog("IlJournalierServer: shutting down", err.Error())
	} else {
		utils.AlwaysLog("IlJournalierServer: shutting down")
	}
}
