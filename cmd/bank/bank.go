package main

import (
	"github.com/DaniilOr/webGo/cmd/bank/app"
	"github.com/DaniilOr/webGo/pkg/CardGiverService"
	"log"
	"net"
	"net/http"
	"os"
)

const defaultPort = "8888"
const defaultHost = "0.0.0.0"

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}

	log.Println(host)
	log.Println(port)

	if err := execute(net.JoinHostPort(host, port)); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(addr string) (err error) {
	cardSvc := CardGiverService.CreateService()
	mux := http.NewServeMux()
	application := app.NewServer(cardSvc, mux)
	application.Init()

	server := &http.Server{
		Addr: addr,
		Handler: application,
	}
	return server.ListenAndServe()
}


