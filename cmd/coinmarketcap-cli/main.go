package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/vlad-golang/coinmarketcap-cli/internal/repo/point_repo_sql"
	"github.com/vlad-golang/coinmarketcap-cli/internal/services/listing_service_impl"
	"github.com/vlad-golang/coinmarketcap-cli/pkg/go-coinmarketcap"
)

func start() error {
	ctx := context.Background()
	client := go_coinmarketcap.NewMarketCapClient()
	pointRepoSql := point_repo_sql.NewPointRepoSql(nil)
	cryptocurrencyServiceImpl := listing_service_impl.NewCryptocurrencyServiceImpl(client, pointRepoSql)
	listing, err := cryptocurrencyServiceImpl.Listing(ctx, listing_service_impl.ListingRequest{
		Limit:  100,
		SortBy: listing_service_impl.CurrentAveragePriceDifference,
	})
	if err != nil {
		return fmt.Errorf("cryptocurrency service impl: %w", err)
	}
	tableWriter := table.NewWriter()
	tableWriter.SetStyle(table.StyleLight)
	tableWriter.AppendHeader(table.Row{
		"#",
		"Name",
		"Price $",
		"Avg $",
		"Min $",
		"Max $",
		"1h %",
		"All %",
		"Created",
		"All chart",
	})
	for _, coin := range listing.Coins {
		plot := asciigraph.Plot(
			coin.Prices,
			asciigraph.Height(2),
			asciigraph.Width(10),
			asciigraph.SeriesColors(asciigraph.Green, asciigraph.Red),
		)

		tableWriter.AppendRow(table.Row{
			coin.Rank,
			coin.Name,
			fmt.Sprintf("%.2f", coin.Price),
			fmt.Sprintf("%.2f", coin.AveragePrice),
			fmt.Sprintf("%.2f", coin.MinPrice),
			fmt.Sprintf("%.2f", coin.MaxPrice),
			fmt.Sprintf("%.2f", coin.PercentChange1h),
			fmt.Sprintf("%.2f", coin.PercentageDifference),
			coin.Created.Format(time.DateOnly),
			plot,
		})

	}
	fmt.Println(tableWriter.Render())

	return nil
}

func main() {
	err := start()
	if err != nil {
		log.Fatal(err)
	}
}
