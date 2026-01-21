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

type customerService struct {
	customerRepository domain.CustomerRepository
}

func NewCustomer(customerRepository domain.CustomerRepository) domain.CustomerService {
	return &customerService{
		customerRepository: customerRepository,
	}
}

func (c customerService) Index(ctx context.Context) ([]dto.CustomerData, error) {
	customers, err := c.customerRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var customerData []dto.CustomerData
	for _, v := range customers {
		customerData = append(customerData, dto.CustomerData{
			ID:   v.ID,
			Code: v.Code,
			Name: v.Name,
		})
	}

	return customerData, nil
}

func (c customerService) Create(ctx context.Context, req dto.CreateRequestCustomers) error {
	customer := domain.Customer{
		Code:      req.Code,
		Name:      req.Name,
		CreatedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}

	return c.customerRepository.Save(ctx, &customer)
}

func (c customerService) Update(ctx context.Context, req dto.CreateCustomerUpdate) error {
	_, err := uuid.Parse(req.ID)
	if err != nil {
		return errors.New("id tidak valid")
	}

	exist, err := c.customerRepository.GetByID(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("id customer tidak ditemukan")
		}
	}

	exist.Code = req.Code
	exist.Name = req.Name
	exist.UpdatedAt = sql.NullTime{Valid: true, Time: time.Now()}

	return c.customerRepository.Update(ctx, &exist)
}

func (c customerService) Delete(ctx context.Context, id string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id tidak valid")
	}
	data, err := c.customerRepository.GetByID(ctx, id)
	if data.ID == "" {
		return errors.New("id tidak di temukan")
	}

	return c.customerRepository.Delete(ctx, id)
}

func (c customerService) Show(ctx context.Context, id string) (dto.CustomerData, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return dto.CustomerData{}, errors.New("id tidak valid")
	}

	data, err := c.customerRepository.GetByID(ctx, id)
	if err != nil {
		return dto.CustomerData{}, err
	}

	if data.ID == "" {
		return dto.CustomerData{}, errors.New("data tidak di temukan")
	}

	return dto.CustomerData{
		ID:   data.ID,
		Code: data.Code,
		Name: data.Name,
	}, nil
}
