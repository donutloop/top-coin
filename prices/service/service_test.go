package service

import (
	"context"
	"github.com/donutloop/top-coin/prices/client"
	"github.com/donutloop/top-coin/prices/proto"
	"net/http"
	"os"
	"testing"
)

func TestPrices_GetPrices(t *testing.T) {

	apiKey := os.Getenv("COIN_MARKET_CAP_APIKEY")

	httpClient := new(http.Client)
	coinMarketCapClient := client.NewCoinMarketCap(httpClient, "https://pro-api.coinmarketcap.com/v1", apiKey, "USD")

	ranksService := NewPrice(coinMarketCapClient)

	resp, err := ranksService.GetPrices(context.Background(), &proto.GetPricesRequest{
		Symbols: []string{"TNCC", "BTC"},
	})

	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Prices) != 1 {
		t.Fatal("symbols count is bad ", resp.Prices)
	}
}
