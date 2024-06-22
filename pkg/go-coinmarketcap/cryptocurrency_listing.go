package go_coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type CryptoCurrency struct {
	ID                            Coin    `json:"id"`
	Name                          string  `json:"name"`
	Symbol                        string  `json:"symbol"`
	Slug                          string  `json:"slug"`
	CMCRank                       int     `json:"cmcRank"`
	MarketPairCount               int     `json:"marketPairCount"`
	CirculatingSupply             float64 `json:"circulatingSupply"`
	SelfReportedCirculatingSupply float64 `json:"selfReportedCirculatingSupply"`
	TotalSupply                   float64 `json:"totalSupply"`
	MaxSupply                     float64 `json:"maxSupply"`
	ATH                           float64 `json:"ath"`
	ATL                           float64 `json:"atl"`
	High24h                       float64 `json:"high24h"`
	Low24h                        float64 `json:"low24h"`
	IsActive                      int     `json:"isActive"`
	LastUpdated                   string  `json:"lastUpdated"`
	DateAdded                     string  `json:"dateAdded"`
	Quotes                        []Quote `json:"quotes"`
}

type Quote struct {
	Name                     string  `json:"name"`
	Price                    float64 `json:"price"`
	Volume24h                float64 `json:"volume24h"`
	Volume7d                 float64 `json:"volume7d"`
	Volume30d                float64 `json:"volume30d"`
	MarketCap                float64 `json:"marketCap"`
	SelfReportedMarketCap    float64 `json:"selfReportedMarketCap"`
	PercentChange1h          float64 `json:"percentChange1h"`
	PercentChange24h         float64 `json:"percentChange24h"`
	PercentChange7d          float64 `json:"percentChange7d"`
	LastUpdated              string  `json:"lastUpdated"`
	PercentChange30d         float64 `json:"percentChange30d"`
	PercentChange60d         float64 `json:"percentChange60d"`
	PercentChange90d         float64 `json:"percentChange90d"`
	FullyDilutedMarketCap    float64 `json:"fullyDilutedMarketCap"`
	MarketCapByTotalSupply   float64 `json:"marketCapByTotalSupply"`
	Dominance                float64 `json:"dominance"`
	Turnover                 float64 `json:"turnover"`
	YtdPriceChangePercentage float64 `json:"ytdPriceChangePercentage"`
	PercentChange1y          float64 `json:"percentChange1y"`
}
type GetCryptocurrencyListingResponse struct {
	Data struct {
		CryptoCurrencyList []CryptoCurrency `json:"cryptoCurrencyList"`
	} `json:"data"`
}

type GetCryptocurrencyListingRequest struct {
	Limit int
}

func (m *MarketCapClient) CryptocurrencyListing(
	ctx context.Context,
	request GetCryptocurrencyListingRequest,
) (response GetCryptocurrencyListingResponse, err error) {
	req, err := http.NewRequest("GET", "https://api.coinmarketcap.com/data-api/v3/cryptocurrency/listing", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	params := url.Values{}
	params.Add("start", "1")
	params.Add("limit", strconv.Itoa(request.Limit))
	params.Add("sortBy", "market_cap")
	params.Add("sortType", "desc")
	params.Add("convert", "USD")
	params.Add("cryptoType", "all")
	params.Add("tagType", "all")
	params.Add("audited", "false")
	params.Add("aux", "ath,atl,high24h,low24h,num_market_pairs,cmc_rank,date_added,max_supply,circulating_supply,total_supply,volume_7d,volume_30d,self_reported_circulating_supply,self_reported_market_cap")
	req.URL.RawQuery = params.Encode()

	m.setHeaders(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, fmt.Errorf("json decode: %v", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return GetCryptocurrencyListingResponse{}, fmt.Errorf("body close: %v", err)
	}

	return response, nil
}
