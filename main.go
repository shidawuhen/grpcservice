/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"fmt"
	"google.golang.org/grpc"
	pb "grpcservice/helloworld"
	"grpcservice/lib"
)
const(
	port = ":50051"
	)


// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.GreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}


//退出前记录错误 & 释放相关资源
func beforeExit() {
	lib.EtcdDelete(port)
	//有错时记录错误信息
	if r := recover(); r != nil {
		tmp := "Panic err : " + r.(string)
		fmt.Println("panic:", tmp)
	}
}

//接收中断
func interrupt() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sign, ok := <-signals
	if ok {
		lib.EtcdDelete(port)
		tmp := "OS Signal received: " + sign.String()
		fmt.Println(tmp)
		os.Exit(0)
	}
}

func main() {
	defer beforeExit()
	go interrupt()
	lib.EtcdPut(port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
