package service

import (
	"context"
	"github.com/donutloop/top-coin/ranks/client"
	"github.com/donutloop/top-coin/ranks/proto"
	"net/http"
	"os"
	"testing"
)

func TestRanks_GetRanks(t *testing.T) {

	apiKey := os.Getenv("CRYPTO_COMPARE_APIKEY")

	httpClient := new(http.Client)
	coinMarketCapClient := client.NewCryptoCompareClient(httpClient, "https://min-api.cryptocompare.com/data/", apiKey, "USD")

	ranksService := NewRanks(coinMarketCapClient)

	resp, err := ranksService.GetRanks(context.Background(), &proto.GetRanksRequest{
		Limit: 200,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Symbols) != 10 {
		t.Fatal("symbols count is bad ", resp.Symbols)
	}
}
