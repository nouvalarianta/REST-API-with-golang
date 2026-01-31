package repository

import (
	"context"
	"database/sql"
	"golang-yt/domain"

	"github.com/doug-martin/goqu/v9"
)

type bookStockRepository struct {
	db *goqu.Database
}

func NewBookStock(con *sql.DB) domain.BookStockRepository {
	return &bookStockRepository{
		db: goqu.New("postgres", con),
	}
}

// GetByBookAndCode implements [domain.BookStockRepository].
func (b *bookStockRepository) GetByBookAndCode(ctx context.Context, id string, code string) (result domain.BookStock, err error) {
	dataset := b.db.From("book_stocks").Where(
		goqu.C("book_id").Eq(id),
		goqu.C("code").Eq(code),
	)
	_, err = dataset.ScanStructContext(ctx, &result)
	return
}

// GetByBookID implements [domain.BookStockRepository].
func (b *bookStockRepository) GetByBookID(ctx context.Context, id string) (result []domain.BookStock, err error) {
	dataset := b.db.From("book_stocks").Where(goqu.C("book_id").Eq(id))
	err = dataset.ScanStructsContext(ctx, &result)
	return
}

// Save implements [domain.BookStockRepository].
func (b *bookStockRepository) Save(ctx context.Context, data []domain.BookStock) error {
	executor := b.db.Insert("book_stocks").Rows(data).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Update implements [domain.BookStockRepository].
func (b *bookStockRepository) Update(ctx context.Context, stock *domain.BookStock) error {
	_, err := b.db.Update("book_stocks").Where(goqu.C("code").Eq(stock.Code)).Set(stock).Executor().ExecContext(ctx)
	// _, err := executor.ExecContext(ctx)
	return err
}

// DeleteBookById implements [domain.BookStockRepository].
func (b *bookStockRepository) DeleteBookById(ctx context.Context, id string) error {
	_, err := b.db.Delete("book_stocks").Where(goqu.C("book_id").Eq(id)).Executor().ExecContext(ctx)
	return err
}

// DeleteByCodes implements [domain.BookStockRepository].
func (b *bookStockRepository) DeleteByCodes(ctx context.Context, codes []string) error {
	_, err := b.db.Delete("book_stocks").Where(goqu.C("code").In(codes)).Executor().ExecContext(ctx)
	return err
}
