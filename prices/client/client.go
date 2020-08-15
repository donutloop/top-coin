package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type (
	MarketQuotesResponse struct {
		Data   map[string]Record `json:"data,omitempty"`
		Status Status            `json:"status"`
	}

	Record struct {
		ID     int                 `json:"id"`
		Symbol string              `json:"symbol"`
		Quote  map[string]Currency `json:"quote"`
	}

	Currency struct {
		Price float64 `json:"price"`
	}

	Status struct {
		ErrorMessage string `json:"error_message,omitempty"`
		ErrorCode    int    `json:"error_code"`
		CreditCount  int    `json:"credit_count"`
	}
)

type CoinMarketCap struct {
	HttpClient   *http.Client
	BaseURL      string
	ApiKey       string
	BaseCurrency string
}

// NewCoinMarketCapClient client for coinmarketcap.com
func NewCoinMarketCap(httpClient *http.Client, baseURL string, APIKey string, BaseCurrency string) CoinMarketCap {
	return CoinMarketCap{
		BaseURL:      baseURL,
		ApiKey:       APIKey,
		BaseCurrency: BaseCurrency,
		HttpClient:   httpClient,
	}
}

func (c CoinMarketCap) GetMarketQuotes(ctx context.Context, names []string) (*MarketQuotesResponse, error) {
	if len(names) == 0 {
		return nil, nil
	}

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse url %v, %w", c.BaseURL, err)
	}
	q := u.Query()
	q.Add("convert", c.BaseCurrency)
	q.Add("symbol", strings.Join(names, ","))
	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, "/cryptocurrency/quotes/latest")
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not init new request %w", err)
	}

	req.Header.Add("X-CMC_PRO_API_KEY", c.ApiKey)
	req.Header.Add("content-type", "application/json")
	
	req = req.WithContext(ctx)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request %w", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body %w", err)
	}

	marketQuotesResponse := new(MarketQuotesResponse)
	if err := json.Unmarshal(respBody, marketQuotesResponse); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body %w", err)
	}

	return marketQuotesResponse, nil
}
