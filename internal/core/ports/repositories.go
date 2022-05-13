package ports

import "github.com/omise/go-tamboon/internal/core/domain"

type DonateRepository interface {
	Get(customerID string) (*domain.Donator, error)
	Scan() ([]domain.Donator, error)
	Create(*domain.Donator) error
	Delete(customerID string) error
	DeleteAll() error
}
