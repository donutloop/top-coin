package main

import (
	"context"
	"fmt"
	ranksproto "github.com/donutloop/top-coin/ranks/proto"
	"github.com/donutloop/top-coin/topcoins/proto"
	"github.com/donutloop/top-coin/topcoins/service"

	pricesproto "github.com/donutloop/top-coin/prices/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

func main() {

	addr := os.Getenv("TOPCOINS_ADDR")
	proxyAddr := os.Getenv("TOPCOINS_PROXY_ADDR")
	pricesAddr := os.Getenv("PRICES_ADDR")
	ranksAddr := os.Getenv("RANKS_ADDR")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to listen on %v", addr))
	}

	ranksGrpcOpts := make([]grpc.DialOption, 0)
	ranksGrpcOpts = append(ranksGrpcOpts, grpc.WithInsecure())

	ranksConn, err := grpc.Dial(ranksAddr, ranksGrpcOpts...)
	if err != nil {
		logrus.Fatalf("could not connect, err: [%v]", err)
	}
	defer func() {
		if err := ranksConn.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	ranksClient := ranksproto.NewRanksServiceClient(ranksConn)

	pricesGrpcOpts := make([]grpc.DialOption, 0)
	pricesGrpcOpts = append(pricesGrpcOpts, grpc.WithInsecure())

	pricesConn, err := grpc.Dial(pricesAddr, pricesGrpcOpts...)
	if err != nil {
		logrus.Fatalf("could not connect, err: [%v]", err)
	}
	defer func() {
		if err := pricesConn.Close(); err != nil {
			logrus.Error(err)
		}
	}()

	pricesClient := pricesproto.NewPricesServiceClient(pricesConn)

	server := grpc.NewServer()
	reflection.Register(server)

	options := make([]grpc.DialOption, 0)
	options = append(options, grpc.WithInsecure())

	topcoinsService := service.NewTopcoins(ranksClient, pricesClient)

	proto.RegisterTopCoinsServiceServer(server, topcoinsService)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("could not start grpc server, err: [%v]\n", err)
			wg.Done()
		}
		wg.Done()
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	wg.Add(1)
	go func() {

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		if err := proto.RegisterTopCoinsServiceHandlerFromEndpoint(ctx, mux, addr, opts); err != nil {
			logrus.Errorf("could not init reverse proxy for service, err: [%v]", err)
		}

		s := &http.Server{
			Addr:    proxyAddr,
			Handler: mux,
		}

		logrus.Infof("proxy server is listing on %s", proxyAddr)
		if err := s.ListenAndServe(); err != nil {
			logrus.Errorf("could listen and serve on http %s, err: [%v]", proxyAddr, err)
		}
		wg.Done()
	}()

	wg.Wait()
}
