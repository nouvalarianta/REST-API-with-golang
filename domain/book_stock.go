package domain

import (
	"context"
	"database/sql"
	"golang-yt/dto"
)

const (
	BookStockStatusAvailable = "AVAILABLE"
	BookStockStatusBorrowed = "BORROWED"
)

type BookStock struct {
	Code string `db:"code"`
	BookID string `db:"book_id"`
	Status string `db:"status"`
	BorrowerID sql.NullString `db:"borrower_id"`
	BorrowedAt sql.NullTime `db:"borrowed_at"`
}

type BookStockRepository interface {
	GetByBookID(ctx context.Context, id string) ([]BookStock, error)
	GetByBookAndCode(ctx context.Context, id string, code string) (BookStock, error)
	Save(ctx context.Context, data []BookStock) error
	Update(ctx context.Context, stock * BookStock) error
	DeleteBookById(ctx context.Context, id string) error
	DeleteByCodes(ctx context.Context, codes []string) error
}

type BookStockService interface {
	Create(ctx context.Context,req dto.CreateBookStockRequest) error
	Delete(ctx context.Context, req dto.DeleteBookStockRequest) error
}