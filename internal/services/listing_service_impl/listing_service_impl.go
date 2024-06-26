package listing_service_impl

import (
	"context"
	"fmt"
	"sort"

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
	Name         string    // Name название
	Price        float64   // Price текущая цена
	Prices       []float64 // Prices история цен
	AveragePrice float64   // AveragePrice средняя цена
	Rank         int
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

			var sum float64
			points := chart.Data.Points
			prices := make([]float64, 0, len(points))
			pointsModels := make([]poin_repo.Point, 0, len(points))
			for _, point := range points {
				price := point.DataPoint.V[0]

				prices = append(prices, price)
				sum += price

				pointsModels = append(pointsModels, poin_repo.Point{
					CoinID:    int64(coin.ID),
					Price:     price,
					TimeStamp: point.Timestamp,
				})
			}

			avg := sum / float64(len(points))
			coins[i] = Coin{coin.Name, usdQuote.Price, prices, avg, coin.CMCRank}
			return nil
		})
	}
	if err = wg.Wait(); err != nil {
		return ListingResponse{}, fmt.Errorf("wg wait: %w", err)
	}

	switch req.SortBy {
	case CurrentAveragePriceDifference:
		sort.Slice(coins, func(i, j int) bool {
			return coins[i].Price-coins[i].AveragePrice < coins[j].Price-coins[j].AveragePrice
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
