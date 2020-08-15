package service

import (
	"context"
	"github.com/donutloop/top-coin/ranks/client"
	"github.com/donutloop/top-coin/ranks/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRanks(CryptoComparer client.CryptoCompare) *Ranks {
	return &Ranks{
		CryptoComparer: CryptoComparer,
	}
}

type Ranks struct {
	CryptoComparer client.CryptoCompare
}

func (r *Ranks) GetRanks(ctx context.Context, req *proto.GetRanksRequest) (*proto.GetRanksResponse, error) {

	// todo(marcel) fetch concurrently all pages

	var symbols []string
	var page int
	var limit int
	if req.Limit > 100 {
		limit = 100
	} else {
		limit = int(req.Limit)
	}

	for req.Limit > 0 {
		if req.Limit < 100 {
			limit = int(req.Limit)
		}

		resp, err := r.CryptoComparer.TopTotalMktCapFull(ctx, limit, page)
		if err != nil {
			logrus.Error(err)
			return nil, status.Errorf(codes.Internal, "could not fetch ranks")
		}

		if len(resp) == 0 {
			break
		}

		for _, coin := range resp {
			symbols = append(symbols, coin.CoinInfo.Name)
		}

		req.Limit = req.Limit - 100
		page++
	}

	getRanksResponse := new(proto.GetRanksResponse)
	getRanksResponse.Symbols = symbols

	return getRanksResponse, nil
}
