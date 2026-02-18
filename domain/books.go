package domain

import (
	"context"
	"database/sql"
	"golang-yt/dto"
)

type Book struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	CoverId     sql.NullString `db:"cover_id"`
	Description string         `db:"description"`
	Isbn        string         `db:"isbn"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	DeletedAt   sql.NullTime   `db:"deleted_at"`
}

type BookRepository interface {
	GetAll(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id string) (Book, error)
	GetByIDs(ctx context.Context, ids []string) ([]Book, error)
	Save(ctx context.Context, book *Book) error
	Update(ctx context.Context, book *Book) error
	Delete(ctx context.Context, id string) error
}

type BookService interface {
	GetAll(ctx context.Context) ([]dto.BookData, error)
	GetByID(ctx context.Context, id string) (dto.BookShowData, error)
	Create(ctx context.Context, req dto.CreateBookRequest) error
	Update(ctx context.Context, req dto.UpdateBookRequest) error
	Delete(ctx context.Context, id string) error
}
