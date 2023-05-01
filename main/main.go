package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	r "github.com/maximotejeda/auth-service/external/v0.0.1/router"
	"github.com/nats-io/nats.go"
)

var (
	port     = os.Getenv("SERVERPORT")
	addr     = os.Getenv("SERVERADDR")
	natsName = os.Getenv("NATSNAME")
)

func main() {

	router := r.NewRouter()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//queue communication between replicas to rotate keys each certain ammount of time
	nc, err := nats.Connect("nats://" + natsName + ":4222/")
	if err != nil {
		log.Print(err)
	}
	ch := make(chan *nats.Msg, 64)
	sub, err := nc.ChanSubscribe("rotate", ch)
	if err != nil {
		log.Print(err)
	}
	defer sub.Unsubscribe()

	r.AuthAddRoutes(router)
	r.UserAddRoutes(router)
	r.AdminAddRoutes(router)

	srv := &http.Server{
		Addr:    addr + ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	for {
		select {
		case <-ch:
			r.J.ReadFromDisk()
		case <-time.After(time.Hour * 24):
			r.J.Renew()
			ch <- nats.NewMsg("renew")
		case <-quit:
			log.Print("Shutting Down server")
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("Server Shutdown:", err)
			}
		}
	}

}
