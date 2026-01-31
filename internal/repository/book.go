package repository

import (
	"context"
	"database/sql"
	"golang-yt/domain"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type bookRepository struct {
	db *goqu.Database
}


func NewBook(con *sql.DB) domain.BookRepository {
	return &bookRepository{
		db: goqu.New("postgres", con),
	}
}

// GetAll implements [domain.BookRepository].
func (b *bookRepository) GetAll(ctx context.Context) (books []domain.Book,err error) {
	dataset := b.db.From("books").Where(goqu.C("deleted_at").IsNull())
	// err = b.db.From("books").Where(goqu.C("deleted_at").IsNull()).ScanStructsContext(ctx, books)
	err = dataset.ScanStructsContext(ctx, &books)
	return
}

// GetByID implements [domain.BookRepository].
func (b *bookRepository) GetByID(ctx context.Context, id string) (book domain.Book,err error) {
	dataset := b.db.From("books").
		Where(
			goqu.C("id").Eq(id),
			goqu.C("deleted_at").IsNull(),
		)
	_, err = dataset.ScanStructContext(ctx, &book)
	return 
}

// Save implements [domain.BookRepository].
func (b *bookRepository) Save(ctx context.Context, book *domain.Book) error {
	executor := b.db.Insert("books").Rows(book).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Update implements [domain.BookRepository].
func (b *bookRepository) Update(ctx context.Context, book *domain.Book) error {
	executor := b.db.Update("books").
		Where(goqu.C("id").Eq(book.ID)).
		Set(book).
		Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

// Delete implements [domain.BookRepository].
func (b *bookRepository) Delete(ctx context.Context, id string) error {
	executor := b.db.Update("books").
		Where(goqu.C("id").Eq(id)).
		Set(goqu.Record{"deleted_at": sql.NullTime{Valid: true, Time: time.Now()}}).
		Executor()
	_, err := executor.ExecContext(ctx)
	return err
} 