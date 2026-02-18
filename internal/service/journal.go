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

type journalService struct {
	journalRepository   domain.JournalRepository
	bookRepository      domain.BookRepository
	bookStockRepository domain.BookStockRepository
	customerRepository  domain.CustomerRepository
	chargeRepository    domain.ChargeRepository
}

func NewJournal(journalRepository domain.JournalRepository, bookRepository domain.BookRepository, bookStockRepository domain.BookStockRepository, customerRepository domain.CustomerRepository, chargeRepository domain.ChargeRepository) domain.JournalService {
	return &journalService{
		journalRepository:   journalRepository,
		bookRepository:      bookRepository,
		bookStockRepository: bookStockRepository,
		customerRepository:  customerRepository,
		chargeRepository:    chargeRepository,
	}
}

// Index implements [domain.JournalService].
func (j *journalService) Index(ctx context.Context, se domain.JournalSearch) ([]dto.JournalData, error) {
	journals, err := j.journalRepository.Find(ctx, se)
	if err != nil {
		return nil, err
	}

	customerId := make([]string, 0)
	bookId := make([]string, 0)

	for _, v := range journals {
		customerId = append(customerId, v.CustomerId)
		bookId = append(bookId, v.BookId)
	}

	customers := make(map[string]domain.Customer)
	if len(customerId) > 0 {
		customersDb, _ := j.customerRepository.GetByIDs(ctx, customerId)
		for _, v := range customersDb {
			customers[v.ID] = v
		}
	}
	books := make(map[string]domain.Book)
	if len(bookId) > 0 {
		bookDb, _ := j.bookRepository.GetByIDs(ctx, bookId)
		for _, v := range bookDb {
			books[v.ID] = v
		}
	}

	result := make([]dto.JournalData, 0)
	for _, v := range journals {
		book := dto.BookData{}
		if v2, e := books[v.BookId]; e {
			book = dto.BookData{
				ID:          v2.ID,
				Isbn:        v2.Isbn,
				Title:       v2.Title,
				Description: v2.Description,
			}
		}

		customer := dto.CustomerData{}
		if v2, e := customers[v.CustomerId]; e {
			customer = dto.CustomerData{
				ID:   v2.ID,
				Code: v2.Code,
				Name: v2.Name,
			}
		}

		result = append(result, dto.JournalData{
			Id:         v.Id,
			BookStock:  v.StockCode,
			Book:       book,
			Customer:   customer,
			BorrowedAt: v.BorrowedAt.Time,
			ReturnedAt: v.ReturnedAt.Time,
		})
	}
	return result, nil
}

// Create implements [domain.JournalService].
func (j *journalService) Create(ctx context.Context, req dto.CreateJournalRequest) error {
	book, err := j.bookRepository.GetByID(ctx, req.BookId)
	if err != nil {
		return err
	}

	if book.ID == "" {
		return errors.New("buku tidak di temukan")
	}

	stock, err := j.bookStockRepository.GetByBookAndCode(ctx, book.ID, req.BookStock)
	if err != nil {
		return err
	}

	if stock.Code == "" {
		return errors.New("buku tidak di temukan")
	}

	if stock.Status != domain.BookStockStatusAvailable {
		return errors.New("buku sudah di pinjam")
	}

	journal := domain.Journal{
		Id:         uuid.NewString(),
		BookId:     req.BookId,
		StockCode:  req.BookStock,
		CustomerId: req.CustomerId,
		Status:     domain.JouralStatusInProgress,
		DueAt:      sql.NullTime{Valid: true, Time: time.Now().Add(7 * 24 * time.Hour)},
		BorrowedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	err = j.journalRepository.Save(ctx, &journal)
	if err != nil {
		return err
	}

	stock.Status = domain.BookStockStatusBorrowed
	stock.BorrowedAt = journal.BorrowedAt
	stock.BorrowerID = sql.NullString{Valid: true, String: journal.CustomerId}

	return j.bookStockRepository.Update(ctx, &stock)
}

// Return implements [domain.JournalService].
func (j *journalService) Return(ctx context.Context, req dto.ReturnJournalRequest) error {
	journal, err := j.journalRepository.FindById(ctx, req.JournaID)
	if err != nil {
		return err
	}

	if journal.Id == "" {
		return domain.JounalNotFound
	}

	stock, err := j.bookStockRepository.GetByBookAndCode(ctx, journal.BookId, journal.StockCode)
	if err != nil {
		return err
	}

	if stock.Code != "" {
		stock.Status = domain.BookStockStatusAvailable
		stock.BorrowerID = sql.NullString{Valid: false}
		stock.BorrowedAt = sql.NullTime{Valid: false}
		err = j.bookStockRepository.Update(ctx, &stock)
		if err != nil {
			return err
		}
	}

	journal.Status = domain.JournalStatusCompleted
	journal.ReturnedAt = sql.NullTime{Valid: true, Time: time.Now()}
	err = j.journalRepository.Update(ctx, &journal)
	if err != nil {
		return err
	}

	hoursLate := time.Since(journal.DueAt.Time).Hours()
	if hoursLate > 24 {

		daysLate := int(hoursLate / 24)

		charge := domain.Charge{
			Id:           uuid.NewString(),
			JournalId:    journal.Id,
			DaysLate:     daysLate,
			DailyLateFee: 5000,
			Total:        5000 * daysLate,
			UserId:       req.UserId,
			CreatedAt:    sql.NullTime{Valid: true, Time: time.Now()},
		}
		err = j.chargeRepository.Save(ctx, &charge)
	}
	return err
}
