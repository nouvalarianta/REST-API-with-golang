package domain

import (
	"context"
	"database/sql"
	"golang-yt/dto"
)

type Customer struct {
	ID        string       `db:"id"`
	Code      string       `db:"code"`
	Name      string       `db:"name"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type CustomerRepository interface {
	GetAll(ctx context.Context) ([]Customer, error)
	GetByID(ctx context.Context, id string) (Customer, error)
	Save(ctx context.Context, c *Customer) error
	Update(ctx context.Context, c *Customer) error
	Delete(ctx context.Context, id string) error
}

type CustomerService interface {
	Index(ctx context.Context) ([]dto.CustomerData, error)
	Create(ctx context.Context, req dto.CreateRequestCustomers) error
	Update(ctx context.Context, req dto.CreateCustomerUpdate) error
	Delete(ctx context.Context, id string) error
	Show(ctx context.Context, id string) (dto.CustomerData, error)
}
