package main

import (
	"context"
	"go-fib/app/config"
	"go-fib/app/internal"
	"go-fib/app/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	conf := config.NewConfig("config.yaml")
	wg := sync.WaitGroup{}
	wg.Add(1)
	g := grpc.NewServer()
	s := internal.NewServer(conf)
	go func() {
		sigc := make(chan os.Signal)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		<-sigc
		g.Stop()
		err := s.Shutdown()
		if err == context.DeadlineExceeded {
			log.Println("shutdown: halted active connections")
		}

		wg.Done()
	}()

	go func() {
		srv := &proto.GRPCServer{Api: s}

		proto.RegisterReverseServer(g, srv)

		l, er := net.Listen("tcp", ":"+conf.GrpcPort)
		if er != nil {
			log.Fatal(er)
		}
		log.Printf("server listening at %v", l.Addr())
		if err := g.Serve(l); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	err := s.Start()
	if err != http.ErrServerClosed {
		log.Fatalln(err)
	}

	wg.Wait()
}
