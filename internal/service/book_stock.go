package service

import (
	"context"
	"errors"
	"golang-yt/domain"
	"golang-yt/dto"
)

type bookStockService struct {
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
}

func NewBookStock(bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository) domain.BookStockService {
	return &bookStockService{
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
	}
}

// Create implements [domain.BookStockService].
func (b *bookStockService) Create(ctx context.Context, req dto.CreateBookStockRequest) error {
	book, err := b.bookRepository.GetByID(ctx, req.BookId)
	if err != nil {
		return err
	}

	if book.ID == "" {
		return errors.New("buku tidak ditemukan")
	}

	stock := make([]domain.BookStock, 0)
	for _, v := range req.Codes {
		stock = append(stock, domain.BookStock{
			Code:   v,
			BookID: req.BookId,
			Status: domain.BookStockStatusAvailable,
		})
	}

	return b.bookStockRepository.Save(ctx, stock)
}

// Delete implements [domain.BookStockService].
func (b *bookStockService) Delete(ctx context.Context, req dto.DeleteBookStockRequest) error {
	return b.bookStockRepository.DeleteByCodes(ctx, req.Codes)
}
