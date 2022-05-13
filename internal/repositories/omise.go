package repositories

import (
	"time"

	"github.com/omise/go-tamboon/internal/core/domain"
	"github.com/omise/go-tamboon/pkg/payment"
)

type omise struct {
	Payment payment.Omise
}

func NewOmiseRepository(public string, secret string, currency string) *omise {

	payment := payment.New(public, secret, currency)

	return &omise{
		Payment: *payment,
	}
}

func (o *omise) Get(customerID string) (*domain.Donator, error) {

	customer, err := o.Payment.RetrieveCustomer(customerID)
	if err != nil {
		return nil, err
	}

	card := customer.Cards.Find(customer.DefaultCard)

	donator := domain.Donator{
		Email:             customer.Email,
		CustomerID:        customer.ID,
		CardID:            card.ID,
		Number:            "",
		SecurityCode:      "",
		Country:           card.Country,
		City:              card.City,
		Bank:              card.Bank,
		PostalCode:        card.PostalCode,
		Financing:         card.Financing,
		LastDigits:        card.LastDigits,
		Brand:             card.Brand,
		ExpirationMonth:   card.ExpirationMonth,
		ExpirationYear:    card.ExpirationYear,
		Fingerprint:       card.Financing,
		Name:              card.Name,
		SecurityCodeCheck: card.SecurityCodeCheck,
	}

	return &donator, nil
}

func (o *omise) Scan() ([]domain.Donator, error) {
	customers, err := o.Payment.RetrieveAllCustomer()
	if err != nil {
		return nil, err
	}

	donators := []domain.Donator{}

	for _, customer := range customers.Data {
		card := customer.Cards.Find(customer.DefaultCard)

		donator := domain.Donator{
			Email:             customer.Email,
			CustomerID:        customer.ID,
			CardID:            card.ID,
			Number:            "",
			SecurityCode:      "",
			Country:           card.Country,
			City:              card.City,
			Bank:              card.Bank,
			PostalCode:        card.PostalCode,
			Financing:         card.Financing,
			LastDigits:        card.LastDigits,
			Brand:             card.Brand,
			ExpirationMonth:   card.ExpirationMonth,
			ExpirationYear:    card.ExpirationYear,
			Fingerprint:       card.Financing,
			Name:              card.Name,
			SecurityCodeCheck: card.SecurityCodeCheck,
		}
		donators = append(donators, donator)
	}

	return donators, nil
}

func (o *omise) Delete(customerID string) error {
	_, err := o.Payment.DeleteCustomer(customerID)
	if err != nil {
		return err
	}

	return nil
}

func (o *omise) DeleteAll() error {
	err := o.Payment.DeleteAllCustomer()
	if err != nil {
		return err
	}

	return nil
}

func (o *omise) Create(donator *domain.Donator) error {
	for {
		customer, err := o.Payment.CreateCustomer(payment.CardDetail{
			Name:            donator.Name,
			Number:          donator.Number,
			City:            donator.City,
			ExpirationMonth: donator.ExpirationMonth,
			ExpirationYear:  donator.ExpirationYear,
			PostalCode:      donator.PostalCode,
			SecurityCode:    donator.SecurityCode,
		})
		if err != nil {
			if err.Error() == "(429/too_many_requests) API rate limit has been exceeded" {
				time.Sleep(5 * time.Second)
				continue
			} else {
				return err
			}
		}
		donator.CustomerID = customer.ID
		break
	}

	// retry when 429/too_many_requests
	for {
		charge, err := o.Payment.CreateCharge(donator.CustomerID, donator.Amount)
		if err != nil {
			if err.Error() == "(429/too_many_requests) API rate limit has been exceeded" {
				time.Sleep(5 * time.Second)
				continue
			} else {
				return err
			}
		}
		donator.ChargeID = charge.ID
		break
	}

	return nil
}
