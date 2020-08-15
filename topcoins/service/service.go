package service

import (
	"context"
	pricesproto "github.com/donutloop/top-coin/prices/proto"
	ranksproto "github.com/donutloop/top-coin/ranks/proto"
	"github.com/donutloop/top-coin/topcoins/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewTopcoins(ranksServiceClient ranksproto.RanksServiceClient, pricesServiceClient pricesproto.PricesServiceClient) *Topcoins {
	return &Topcoins{
		RanksServiceClient:  ranksServiceClient,
		PricesServiceClient: pricesServiceClient,
	}
}

type Topcoins struct {
	RanksServiceClient  ranksproto.RanksServiceClient
	PricesServiceClient pricesproto.PricesServiceClient
}

func (r *Topcoins) GetTopcoins(ctx context.Context, req *proto.GetTopcoinsRequest) (*proto.GetTopcoinsResponse, error) {

	getRanksRequest := &ranksproto.GetRanksRequest{
		Limit: req.Limit,
	}

	getRanksResponse, err := r.RanksServiceClient.GetRanks(ctx, getRanksRequest)
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "could not fetch ranks")
	}

	getPricesRequest := &pricesproto.GetPricesRequest{
		Symbols: getRanksResponse.Symbols,
	}

	getPricesResponse, err := r.PricesServiceClient.GetPrices(ctx, getPricesRequest)
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "could not fetch prices")
	}

	var coins []*proto.Coin
	for i, symbol := range getRanksResponse.Symbols {
		price, ok := getPricesResponse.Prices[symbol]
		if !ok {
			price = -1
		}
		coins = append(coins, &proto.Coin{
			Price:  price,
			Symbol: symbol,
			Rank:   int64(i + 1),
		})
	}

	getTopcoinsResponse := new(proto.GetTopcoinsResponse)
	getTopcoinsResponse.Coin = coins

	return getTopcoinsResponse, nil
}
