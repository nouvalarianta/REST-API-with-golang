package repository

import (
	"context"
	"database/sql"
	"golang-yt/domain"

	"github.com/doug-martin/goqu/v9"
)

type chargeRepository struct {
	db *goqu.Database
}

func NewCharge(con *sql.DB) domain.ChargeRepository {
	return &chargeRepository{
		db: goqu.New("postgres", con),
	}
}


// Save implements [domain.ChargeRepository].
func (c *chargeRepository) Save(ctx context.Context, charge *domain.Charge) error {
	_, err := c.db.Insert("charges").Rows(charge).Executor().ExecContext(ctx)
	return  err
}