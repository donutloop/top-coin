package service

import (
	"context"
	"github.com/donutloop/top-coin/prices/client"
	"github.com/donutloop/top-coin/prices/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type Price struct {
	coinMarketCapClient client.CoinMarketCap
}

func NewPrice(CoinMarketCapClient client.CoinMarketCap) *Price {
	return &Price{coinMarketCapClient: CoinMarketCapClient}
}

func (s *Price) GetPrices(ctx context.Context, req *proto.GetPricesRequest) (*proto.GetPricesResponse, error) {

	data, err := s.coinMarketCapClient.GetMarketQuotes(ctx, req.Symbols)
	if err != nil {
		logrus.Error(err)
		return nil, status.Errorf(codes.Internal, "could not fetch prices")
	}

	// todo(marcel) Save all invalid symbols in database or cache
	if data.Status.ErrorCode == 400  {
		if strings.Contains(data.Status.ErrorMessage, "Invalid values for \"symbol\"") || strings.Contains(data.Status.ErrorMessage, "Invalid value for \"symbol\"") {
			filteredSymbols := make([]string, 0)
			for _, symbol := range req.Symbols {
				if !strings.Contains(data.Status.ErrorMessage, symbol) {
					filteredSymbols = append(filteredSymbols, symbol)
				}
			}
			if len(filteredSymbols) == 0 {
				return nil, status.Errorf(codes.Unknown, "Invalid symbol values")
			}

			data, err = s.coinMarketCapClient.GetMarketQuotes(ctx, filteredSymbols)
			if err != nil {
				logrus.Error(err)
				return nil, status.Errorf(codes.Internal, "could not fetch prices")
			}

			if data.Status.ErrorCode != 0 {
				return nil, status.Errorf(codes.Unknown, data.Status.ErrorMessage)
			}
		} else {
			return nil, status.Errorf(codes.Unknown, data.Status.ErrorMessage)
		}
	}

	getPricesResponse := new(proto.GetPricesResponse)
	getPricesResponse.Prices = make(map[string]float64)
	for _, record := range data.Data {
		getPricesResponse.Prices[record.Symbol] = record.Quote[s.coinMarketCapClient.BaseCurrency].Price
	}
	return getPricesResponse, nil
}
