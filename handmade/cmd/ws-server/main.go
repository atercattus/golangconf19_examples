package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	argv struct {
		port int
		help bool
	}

	wsServer *WSServer
)

func init() {
	flag.IntVar(&argv.port, `port`, 8081, `listen port`)
	flag.BoolVar(&argv.help, `h`, false, `show this help`)
	flag.Parse()
}

func main() {
	wsServer = MakeWSServer(fmt.Sprintf(`:%d`, argv.port))
	log.Printf(`Started at port %d`, argv.port)
	go func() {
		if err := wsServer.httpServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	waitSigQuit()
}

func waitSigQuit() {
	sig := make(chan os.Signal, 10)
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGINT)
	<-sig
}
