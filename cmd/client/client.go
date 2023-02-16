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
	"strings"
	"time"

	"github.com/ciena/grpc-hello/internal/pkg/info"
	pb "github.com/ciena/grpc-hello/pkg/apis/hello"
	"github.com/ciena/grpc-hello/pkg/version"

	"google.golang.org/grpc"
)

type config struct {
	Addr         string
	SendInterval time.Duration
	Timeout      time.Duration
	ShowVersion  bool
}

func main() {
	var cfg config
	flag.StringVar(&cfg.Addr,
		"addr", "localhost:8080", "address and port on which to send hello")
	flag.DurationVar(&cfg.SendInterval,
		"interval", 3*time.Second, "how often to send hello")
	flag.DurationVar(&cfg.Timeout,
		"timeout", 3*time.Second, "how long to wait for a response")
	flag.BoolVar(&cfg.ShowVersion,
		"version", false, "display version and exit")
	flag.Parse()

	if cfg.ShowVersion {
		fmt.Println(version.Version().String())
		os.Exit(0)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	hostname, _ := os.Hostname()
	ip, _ := info.MyIP()
	for {
		select {
		case <-sigs:
			log.Println("Interrupt requested .... exit")
			os.Exit(0)
		case <-time.After(cfg.SendInterval):
			parts := strings.Split(cfg.Addr, ":")
			addrs, err := net.LookupHost(parts[0])
			if err != nil {
				log.Printf("unable to resolve '%s': %v", parts[0], err)
			} else {
				log.Printf("server '%s' resolved as %s", parts[0], strings.Join(addrs, ", "))
			}
			conn, err := grpc.Dial(cfg.Addr, grpc.WithInsecure())
			if err != nil {
				log.Printf("did not connect: %v", err)
				break
			}
			c := pb.NewHelloClient(conn)

			ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
			defer cancel()
			sentAt := time.Now().UTC()
			resp, err := c.SayHello(ctx, &pb.HelloRequest{
				SentAt: sentAt.Format("Mon Jan 2 15:04:05.000 MST 2006"),
				SentBy: &pb.ID{
					IP:       ip.String(),
					Hostname: hostname}})
			if err != nil {
				log.Printf("could not greet: %v", err)
			}
			conn.Close()
			delta := time.Now().UTC().Sub(sentAt)
			log.Printf("Greeting response: %v (%v)", resp, delta)
		}
	}
}
