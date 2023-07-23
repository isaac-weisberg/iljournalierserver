package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("IlJournalierServer start!")

	ctx := context.Background()

	diContainer, err := newDIContainer(ctx)

	if err != nil {
		panic(err)
	}

	mainRouter := newMainRouter(diContainer)

	mux := http.NewServeMux()

	mux.Handle("/", &mainRouter)

	listenAddress := ":24610"
	server := http.Server{Addr: listenAddress, Handler: mux}

	go func() {
		fmt.Println("IlJournalierServer will listen", listenAddress)
		err := server.ListenAndServe()
		if err != nil {
			println("IlJournalierServer returned", err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		fmt.Println("IlJournalierServer shutting down", err.Error())
	} else {
		fmt.Println("IlJournalierServer shutting down")
	}
}
