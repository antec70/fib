package main

import (
	"context"
	"fmt"
	"go-fib/app/proto"
	"google.golang.org/grpc"
	"log"
	"time"
)

const address = "localhost:50051"

func main() {
	conn, er := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if er != nil {
		log.Fatal(er)
	}
	defer conn.Close()
	c := proto.NewReverseClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Fib(ctx, &proto.Request{X: 0, Y: 16})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r)
}
