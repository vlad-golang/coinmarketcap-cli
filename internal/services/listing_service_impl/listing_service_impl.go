package listing_service_impl

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/vlad-golang/coinmarketcap-cli/internal/interfaces/repo/poin_repo"
	"github.com/vlad-golang/coinmarketcap-cli/internal/repo/point_repo_sql"
	"github.com/vlad-golang/coinmarketcap-cli/pkg/go-coinmarketcap"
)

//go:generate gonstructor --type=CryptocurrencyServiceImpl --constructorTypes=allArgs --output=./constructor.go
type CryptocurrencyServiceImpl struct {
	coinMarketCapClient *go_coinmarketcap.MarketCapClient
	pointRepo           *point_repo_sql.PointRepoSql
}

type ListingRequestSortBy string

const (
	SortByAveragePrice            ListingRequestSortBy = "average_price"
	CurrentAveragePriceDifference                      = "currentAveragePriceDifference"
)

type ListingRequest struct {
	Limit      int // Limit Количество монет
	SortBy     ListingRequestSortBy
	IsDescSort bool
}

// Coin монета
type Coin struct {
	Name                 string    // Name название
	Price                float64   // Price текущая цена
	Prices               []float64 // Prices история цен
	AveragePrice         float64   // AveragePrice средняя цена
	Rank                 int
	PercentageDifference float64
	PercentChange1h      float64
	PercentChange24h     float64
	PercentChange7d      float64
	MaxPrice             float64
	MinPrice             float64
	Created              time.Time
}

type ListingResponse struct {
	Coins []Coin
}

func (c *CryptocurrencyServiceImpl) Listing(ctx context.Context, req ListingRequest) (ListingResponse, error) {
	listing, err := c.coinMarketCapClient.CryptocurrencyListing(ctx, go_coinmarketcap.GetCryptocurrencyListingRequest{Limit: req.Limit})
	if err != nil {
		return ListingResponse{}, fmt.Errorf("coin market cap client get cryptocurrency listing: %w", err)
	}

	coins := make([]Coin, len(listing.Data.CryptoCurrencyList)) // Создание слайса с нужной длиной для сохранения порядка
	semaphore := make(chan struct{}, 10)                        // Ограничение до 10 параллельных запросов

	wg := errgroup.Group{}
	for i, coin := range listing.Data.CryptoCurrencyList {
		i, coin := i, coin // создание копий переменных для использования в горутине
		usdQuote := coin.Quotes[len(coin.Quotes)-1]

		wg.Go(func() error {
			semaphore <- struct{}{}        // Захват семафора
			defer func() { <-semaphore }() // Освобождение семафора

			chart, err := c.coinMarketCapClient.CryptocurrencyDetailChart(ctx, coin.ID)
			if err != nil {
				return fmt.Errorf("coin market cap client cryptocurrency detail chart: %w", err)
			}

			var sumPrice float64
			var maxPrice float64
			minPrice := math.MaxFloat64
			points := chart.Data.Points
			created := time.Unix(points[0].Timestamp, 0)
			prices := make([]float64, 0, len(points))
			pointsModels := make([]poin_repo.Point, 0, len(points))
			for _, point := range points {
				price := point.DataPoint.V[0]
				if price > maxPrice {
					maxPrice = price
				}
				if price < minPrice {
					minPrice = price
				}

				prices = append(prices, price)
				sumPrice += price

				pointsModels = append(pointsModels, poin_repo.Point{
					CoinID:    int64(coin.ID),
					Price:     price,
					TimeStamp: point.Timestamp,
				})
			}

			avgPrice := sumPrice / float64(len(points))
			difference := usdQuote.Price - avgPrice
			percentageDifference := (difference / avgPrice) * 100
			coins[i] = Coin{coin.Name, usdQuote.Price, prices, avgPrice,
				coin.CMCRank, percentageDifference, usdQuote.PercentChange1h,
				usdQuote.PercentChange24h, usdQuote.PercentChange30d, maxPrice, minPrice, created}
			return nil
		})
	}
	if err = wg.Wait(); err != nil {
		return ListingResponse{}, fmt.Errorf("wg wait: %w", err)
	}

	switch req.SortBy {
	case CurrentAveragePriceDifference:
		sort.Slice(coins, func(i, j int) bool {
			return coins[i].PercentageDifference < coins[j].PercentageDifference
		})
	case SortByAveragePrice:
		sort.Slice(coins, func(i, j int) bool {
			return coins[i].AveragePrice < coins[j].AveragePrice
		})
	}

	return ListingResponse{
		Coins: coins,
	}, nil
}
