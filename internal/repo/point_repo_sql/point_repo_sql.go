package point_repo_sql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vlad-golang/coinmarketcap-cli/internal/interfaces/repo/poin_repo"
)

//go:generate gonstructor --type=PointRepoSql --constructorTypes=allArgs --output=./constructor.go
type PointRepoSql struct {
	db *sql.DB
}

// goverter:converter
// goverter:output:file ./converter/converter.gen.go
// goverter:output:package :converter
// goverter:enum:unknown @error
//
//go:generate goverter gen ./
type WebAdminHttpConverter interface {
	//PointFromRepoToJet([]poin_repo.Point) []model.Points
}

func (p *PointRepoSql) CreateOrUpdate(ctx context.Context, points []poin_repo.Point) error {
	return errors.New("not implemented")
}

func (p *PointRepoSql) All(ctx context.Context) ([]poin_repo.Point, error) {
	return nil, errors.New("not implemented")
}
