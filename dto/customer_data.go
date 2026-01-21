package dto

type CustomerData struct {
	ID string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type CreateRequestCustomers struct{
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type CreateCustomerUpdate struct{
	ID string `json:"-"`
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}