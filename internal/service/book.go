package service

import (
	"context"
	"database/sql"
	"errors"
	"golang-yt/domain"
	"golang-yt/dto"
	"time"

	"github.com/google/uuid"
)

type bookService struct {
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
}

func NewBook(bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository) domain.BookService {
	return &bookService{
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
	}
}

// GetAll implements [domain.BookService].
func (b *bookService) GetAll(ctx context.Context) (data []dto.BookData,err error) {
	result, err := b.bookRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	//kurang paham
	for _, v := range result {
		data = append(data, dto.BookData{
			ID: v.ID,
			Isbn: v.Isbn,
			Title: v.Title,
			Description: v.Description,
		})
	}

	return data, nil
}

// GetByID implements [domain.BookService].
func (b *bookService) GetByID(ctx context.Context, id string) (dto.BookShowData, error) {
	data, err := b.bookRepository.GetByID(ctx, id)
	if err != nil {
		return dto.BookShowData{}, err
	}

	if data.ID == "" {
		return dto.BookShowData{}, errors.New("buku tidak di temukan")
	}

	stocks, err := b.bookStockRepository.GetByBookID(ctx, data.ID, )
	if err != nil {
		return dto.BookShowData{}, err
	}

	stockData := make([]dto.BookStockData, 0)
	for _, v := range stocks {
		stockData = append(stockData, dto.BookStockData{
			Code: v.Code,
			Status: v.Status,
		})
	}
	
	return dto.BookShowData{
		BookData: dto.BookData{
		ID: data.ID,
		Isbn: data.Isbn,
		Title: data.Title,
		Description: data.Description,
		}, 
		Stocks: stockData,
	}, err
}
// Create implements [domain.BookService].
func (b *bookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	book := domain.Book{
		ID: uuid.NewString(),
		Title: req.Title,
		Description: req.Description,
		Isbn: req.Isbn,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	return b.bookRepository.Save(ctx, &book)
}

// Update implements [domain.BookService].
func (b *bookService) Update(ctx context.Context, req dto.UpdateBookRequest) error {
	persisted, err := b.bookRepository.GetByID(ctx, req.Id)
	if err != nil {
		return err
	}

	if persisted.ID == " " {
		return errors.New("buku tidak di temukan")
	}
	 
	persisted.Title = req.Title
	persisted.Isbn = req.Isbn
	persisted.Description = req.Description
	persisted.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}

	return b.bookRepository.Update(ctx, &persisted)
}

// Delete implements [domain.BookService].
func (b *bookService) Delete(ctx context.Context, id string) error {
	persisted, err := b.bookRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if persisted.ID == " " {
		return errors.New("buku tidak di temukan")
	}

	err = b.bookRepository.Delete(ctx, persisted.ID)
	if err != nil {
		return err
	}

	return b.bookStockRepository.DeleteBookById(ctx, persisted.ID)
}