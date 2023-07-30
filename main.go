package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"caroline-weisberg.fun/iljournalierserver/intake"
	"caroline-weisberg.fun/iljournalierserver/utils"
)

func main() {
	var intakeConfig, err = intake.ReadIntakeConfiguration()
	if err != nil {
		utils.AlwaysLog(err.Error())
		return
	}

	utils.AlwaysLog("IlJournalierServer: start!")

	ctx := context.Background()

	diContainer, err := newDIContainer(ctx, intakeConfig)

	if err != nil {
		panic(err)
	}

	mainRouter := newMainRouter(diContainer)

	mux := http.NewServeMux()

	mux.Handle("/", &mainRouter)

	listenAddress := ":24610"
	server := http.Server{Addr: listenAddress, Handler: mux}

	go func() {
		utils.AlwaysLog("IlJournalierServer: will listen", listenAddress)
		err := server.ListenAndServe()
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
