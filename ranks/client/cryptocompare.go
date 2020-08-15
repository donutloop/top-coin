package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type (
	Currency struct {
		Symbol         string  `json:"SYMBOL"`
		Supply         float64 `json:"SUPPLY"`
		Fullname       string  `json:"FULLNAME"`
		Name           string  `json:"NAME"`
		Volume24HourTo float64 `json:"VOLUME24HOURTO"`
	}

	TopListByPairVolumeResponse struct {
		Data           []Currency `json:"Data,omitempty"`
		Type           int        `json:"Type"`
		Response       string     `json:"Response"`
		Message        string     `json:"Message"`
		HasWarning     bool       `json:"HasWarning"`
		ParamWithError string     `json:"ParamWithError,omitempty"`
	}

	ToplistByMarketCapFullDataResponse struct {
		Data    []CoinData `json:"Data,omitempty"`
		Message string     `json:"Message"`
	}

	CoinData struct {
		CoinInfo CoinInfo `json:"CoinInfo"`
	}
	CoinInfo struct {
		Name string `json:"Name"`
	}
)


type CryptoCompare struct {
	HttpClient   *http.Client
	BaseURL      string
	APIkey       string
	BaseCurrency string
}

func NewCryptoCompareClient(httpClient *http.Client, BaseURL string, APIKey string, BaseCurrency string) CryptoCompare {
	return CryptoCompare{
		BaseURL:      BaseURL,
		APIkey:       APIKey,
		BaseCurrency: BaseCurrency,
		HttpClient:   httpClient,
	}
}

// TopTotalMktCapFull get ordered list of currencies https://min-api.cryptocompare.com/documentation?key=Toplists&cat=TopTotalMktCapEndpointFull
func (c CryptoCompare) TopTotalMktCapFull(ctx context.Context, limit, page int) ([]CoinData, error) {

	v := url.Values{}
	v.Add("tsym", c.BaseCurrency)
	v.Add("limit", strconv.Itoa(limit))
	v.Add("page", strconv.Itoa(page))

	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse url %v, %w", c.BaseURL, err)
	}
	u.RawQuery = v.Encode()

	u.Path = path.Join(u.Path, "/top/mktcapfull/")
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not init new request %w", err)
	}

	req.Header.Add("authorization", fmt.Sprintf("Apikey %s", c.APIkey))
	req.Header.Add("content-type", "application/json")

	req = req.WithContext(ctx)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response code is bad %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body %w", err)
	}

	toplistByMarketCapFullDataResponse := new(ToplistByMarketCapFullDataResponse)
	if err := json.Unmarshal(respBody, toplistByMarketCapFullDataResponse); err != nil {
		return nil, fmt.Errorf("could not unmarshal response body %w", err)
	}

	return toplistByMarketCapFullDataResponse.Data, nil
}
