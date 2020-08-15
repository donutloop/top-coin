package main

import (
	"fmt"
	"github.com/donutloop/top-coin/prices/client"
	"github.com/donutloop/top-coin/prices/proto"
	"github.com/donutloop/top-coin/prices/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

func main() {

	addr := os.Getenv("PRICES_ADDR")
	apiKey := os.Getenv("COIN_MARKET_CAP_APIKEY")

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to listen on %v", addr))
	}

	server := grpc.NewServer()
	reflection.Register(server)

	options := make([]grpc.DialOption, 0)
	options = append(options, grpc.WithInsecure())

	httpClient := new(http.Client)
	coinMarketCapClient := client.NewCoinMarketCap(httpClient, "https://pro-api.coinmarketcap.com/v1", apiKey, "USD")

	pricesService := service.NewPrice(coinMarketCapClient)

	proto.RegisterPricesServiceServer(server, pricesService)

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
