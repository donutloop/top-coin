package main

import (
	"fmt"
	"github.com/donutloop/top-coin/ranks/client"
	"github.com/donutloop/top-coin/ranks/proto"
	"github.com/donutloop/top-coin/ranks/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

func main() {

	addr := os.Getenv("RANKS_ADDR")
	apiKey := os.Getenv("CRYPTO_COMPARE_APIKEY")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to listen on %v", addr))
	}

	server := grpc.NewServer()
	reflection.Register(server)

	options := make([]grpc.DialOption, 0)
	options = append(options, grpc.WithInsecure())

	httpClient := new(http.Client)
	coinMarketCapClient := client.NewCryptoCompareClient(httpClient, "https://min-api.cryptocompare.com/data/", apiKey, "USD")

	ranksService := service.NewRanks(coinMarketCapClient)

	proto.RegisterRanksServiceServer(server, ranksService)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalf("could not start grpc server, err: [%v]\n", err)
			wg.Done()
		}
		wg.Done()
	}()

	wg.Wait()
}
