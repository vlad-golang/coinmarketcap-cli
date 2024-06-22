package main

import (
	"context"
	"fmt"
	"log"

	"github.com/guptarohit/asciigraph"

	"github.com/vlad-golang/coinmarketcap-cli/internal/repo/point_repo_sql"
	"github.com/vlad-golang/coinmarketcap-cli/internal/services/listing_service_impl"
	"github.com/vlad-golang/coinmarketcap-cli/pkg/go-coinmarketcap"
)

func start() error {
	ctx := context.Background()
	client := go_coinmarketcap.NewMarketCapClient()
	pointRepoSql := point_repo_sql.NewPointRepoSql(nil)
	cryptocurrencyServiceImpl := listing_service_impl.NewCryptocurrencyServiceImpl(client, pointRepoSql)
	listing, err := cryptocurrencyServiceImpl.Listing(ctx, listing_service_impl.ListingRequest{Limit: 100})
	if err != nil {
		return fmt.Errorf("cryptocurrency service impl: %w", err)
	}

	for _, coin := range listing.Coins {
		avgPoints := make([]float64, 0, len(coin.Prices))
		for range coin.Prices {
			avgPoints = append(avgPoints, coin.AveragePrice)
		}

		plot := asciigraph.PlotMany(
			[][]float64{
				coin.Prices,
				avgPoints,
			},
			asciigraph.Height(10),
			asciigraph.Width(50),
			asciigraph.Caption(coin.Name),
			asciigraph.SeriesColors(asciigraph.Green, asciigraph.Red),
			asciigraph.SeriesLegends("Price", "Average"),
		)
		fmt.Println("Name:", coin.Name)
		fmt.Println("Price:", coin.Price)
		fmt.Println("Average price:", coin.AveragePrice)
		fmt.Println(plot)
		fmt.Println()
	}

	return nil
}

func main() {
	err := start()
	if err != nil {
		log.Fatal(err)
	}
}
