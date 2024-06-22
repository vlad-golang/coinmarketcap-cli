package go_coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

type CryptocurrencyDetailRequest struct {
}

type Range string

const (
	All    Range = "ALL"
	OneDay Range = "1D"
)

type CryptocurrencyDetailChartResponse struct {
	Data struct {
		Points []TimedDataPoint
	}
}
type TimedDataPoint struct {
	Timestamp int64
	DataPoint DataPoint
}

type DataPoint struct {
	V []float64 `json:"v"`
	C []float64 `json:"c"`
}

func (m *MarketCapClient) CryptocurrencyDetailChart(ctx context.Context, coin Coin) (response CryptocurrencyDetailChartResponse, err error) {
	type CryptocurrencyDetailChartHttpResponse struct {
		Data struct {
			Points map[int64]DataPoint `json:"points"`
		} `json:"data"`
	}

	url := fmt.Sprintf("https://api.coinmarketcap.com/data-api/v3/cryptocurrency/detail/chart?id=%d&range=ALL", coin)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	m.setHeaders(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	var httpResponse CryptocurrencyDetailChartHttpResponse
	err = json.NewDecoder(resp.Body).Decode(&httpResponse)
	if err != nil {
		return response, err
	}

	for key, value := range httpResponse.Data.Points {
		response.Data.Points = append(response.Data.Points, TimedDataPoint{
			Timestamp: key,
			DataPoint: DataPoint{
				V: value.V,
				C: value.C,
			},
		})
	}

	sort.Slice(response.Data.Points, func(i, j int) bool {
		return response.Data.Points[i].Timestamp < response.Data.Points[j].Timestamp
	})

	return response, nil
}
