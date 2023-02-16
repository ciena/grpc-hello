/*
Copyright 2021 Ciena Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ciena/grpc-hello/internal/pkg/info"
	pb "github.com/ciena/grpc-hello/pkg/apis/hello"
	"github.com/ciena/grpc-hello/pkg/version"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// How to get ip from ctx ?
	var peerIP string
	md, ok := peer.FromContext(ctx)
	if ok {
		peerIP = md.Addr.String()
	}

	now := time.Now().UTC()
	sentAt, _ := time.Parse("Mon Jan 2 15:04:05.000 MST 2006", req.SentAt)
	log.Printf("Greeting received: %v (%v)", req, now.Sub(sentAt))
	log.Printf("  gRPC Peer IP: %v", peerIP)

	hostname, _ := os.Hostname()
	ip, _ := info.MyIP()

	return &pb.HelloResponse{
		ReceivedAt: time.Now().UTC().Format("Mon Jan 2 15:04:05.000 MST 2006"),
		RespondedBy: &pb.ID{
			IP:       ip.String(),
			Hostname: hostname}}, nil
}

type config struct {
	ListenOn    string
	ShowVersion bool
}

func main() {
	var cfg config
	flag.StringVar(&cfg.ListenOn, "listen", ":8080", "address and port on which to listen")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "show version and exit")
	flag.Parse()

	if cfg.ShowVersion {
		fmt.Println(version.Version().String())
		os.Exit(0)
	}

	lis, err := net.Listen("tcp", cfg.ListenOn)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &server{})
	reflection.Register(s)
	log.Printf("accepting requests on '%s'\n", cfg.ListenOn)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case <-sigs:
		log.Println("Interrupt requested .... exit")
		os.Exit(0)
	}
}
