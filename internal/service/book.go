package service

import (
	"context"
	"database/sql"
	"errors"
	"golang-yt/domain"
	"golang-yt/dto"
	"golang-yt/internal/config"
	"path"
	"time"

	"github.com/google/uuid"
)

type bookService struct {
	cnf                 *config.Config
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
	mediaRepository     domain.MediaRepository
}

func NewBook(cnf *config.Config, bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository, mediaRepository domain.MediaRepository) domain.BookService {
	return &bookService{
		cnf:                 cnf,
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
		mediaRepository:     mediaRepository,
	}
}

// GetAll implements [domain.BookService].
func (b *bookService) GetAll(ctx context.Context) (data []dto.BookData, err error) {
	result, err := b.bookRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	//
	coverId := make([]string, 0)
	for _, v := range result {
		if v.CoverId.Valid {
			coverId = append(coverId, v.CoverId.String)
		}
	}

	covers := make(map[string]string)
	if len(coverId) > 0 {
		coverDb, _ := b.mediaRepository.FindByIds(ctx, coverId)
		for _, v := range coverDb {
			covers[v.Id] = path.Join(b.cnf.Server.Asset, v.Path)
		}
	}

	//
	for _, v := range result {
		var coverUrl string
		if v2, e := covers[v.CoverId.String]; e {
			coverUrl = v2
		}

		data = append(data, dto.BookData{
			ID:          v.ID,
			Isbn:        v.Isbn,
			Title:       v.Title,
			CoverUrl:    coverUrl,
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

	stocks, err := b.bookStockRepository.GetByBookID(ctx, data.ID)
	if err != nil {
		return dto.BookShowData{}, err
	}

	stockData := make([]dto.BookStockData, 0)
	for _, v := range stocks {
		stockData = append(stockData, dto.BookStockData{
			Code:   v.Code,
			Status: v.Status,
		})
	}

	var coverUrl string
	if data.CoverId.Valid {
		cover, _ := b.mediaRepository.FindById(ctx, data.CoverId.String)
		if cover.Path != "" {
			coverUrl = path.Join(b.cnf.Server.Asset, cover.Path) 
		}
	}

	return dto.BookShowData{
		BookData: dto.BookData{
			ID:          data.ID,
			Isbn:        data.Isbn,
			Title:       data.Title,
			CoverUrl: coverUrl,
			Description: data.Description,
		},
		Stocks: stockData,
	}, err
}

// Create implements [domain.BookService].
func (b *bookService) Create(ctx context.Context, req dto.CreateBookRequest) error {
	coverId := sql.NullString{Valid: false, String: req.CoverId}
	if req.CoverId != "" {
		coverId.Valid = true
	}

	book := domain.Book{
		ID:          uuid.NewString(),
		Title:       req.Title,
		CoverId:     coverId,
		Description: req.Description,
		Isbn:        req.Isbn,
		CreatedAt:   sql.NullTime{Valid: true, Time: time.Now()},
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

	coverId := sql.NullString{Valid: false, String: req.CoverId}
	if req.CoverId != "" {
		coverId.Valid = true
	}
	persisted.Title = req.Title
	persisted.Isbn = req.Isbn
	persisted.Description = req.Description
	persisted.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}
	persisted.CoverId = coverId

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
