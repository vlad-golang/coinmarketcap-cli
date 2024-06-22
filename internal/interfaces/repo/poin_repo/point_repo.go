package poin_repo

import "context"

type Point struct {
	ID        int64
	CoinID    int64
	Price     float64
	TimeStamp int64
}

// PointRepo Репозиторий точек для построения графика цены
type PointRepo interface {
	// Create Создает точки
	Create(ctx context.Context) ([]Point, error)
}
