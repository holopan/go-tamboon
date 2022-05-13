package payment

import (
	"time"

	omiseLib "github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type Omise struct {
	Currency string
	client   omiseLib.Client
	secret   string
}

type CardDetail struct {
	Name            string
	Number          string
	City            string
	ExpirationMonth time.Month
	ExpirationYear  int
	PostalCode      string
	SecurityCode    string
}

func New(public string, secret string, currency string) *Omise {

	o := &Omise{}
	client, err := omiseLib.NewClient(public, secret)
	if err != nil {
		panic(err)
	}

	o.client = *client
	o.Currency = currency
	o.secret = secret

	return o
}

func (o *Omise) createCard(cardDetail CardDetail) (*omiseLib.Card, error) {

	card, createToken := &omiseLib.Card{}, &operations.CreateToken{
		Name:            cardDetail.Name,
		Number:          cardDetail.Number,
		ExpirationMonth: cardDetail.ExpirationMonth,
		ExpirationYear:  cardDetail.ExpirationYear,
		City:            cardDetail.City,
		PostalCode:      cardDetail.PostalCode,
		SecurityCode:    cardDetail.SecurityCode,
	}
	if err := o.client.Do(card, createToken); err != nil {

		return nil, err
	}

	return card, nil
}

func (o *Omise) CreateCustomer(cardDetail CardDetail) (*omiseLib.Charge, error) {

	token, err := o.createCard(cardDetail)
	if err != nil {
		return nil, err
	}

	charge, payload := &omiseLib.Charge{}, &operations.CreateCustomer{
		Card: token.ID,
	}
	err = o.client.Do(charge, payload)
	if err != nil {
		return nil, err
	}

	return charge, nil
}

func (o *Omise) RetrieveCustomer(customerID string) (*omiseLib.Customer, error) {

	customer, payload := &omiseLib.Customer{}, &operations.RetrieveCustomer{
		CustomerID: customerID,
	}
	err := o.client.Do(customer, payload)
	if err != nil {
		return nil, err
	}

	return customer, nil

}

func (o *Omise) RetrieveAllCustomer() (*omiseLib.CustomerList, error) {

	customers, list := &omiseLib.CustomerList{}, &operations.ListCustomers{
		List: operations.List{
			From:  time.Now().Add(-1 * time.Hour),
			Limit: 100,
		},
	}
	if err := o.client.Do(customers, list); err != nil {
		return nil, err
	}

	return customers, nil

}

func (o *Omise) CreateCharge(customer string, amount int64) (*omiseLib.Charge, error) {

	charge, payload := &omiseLib.Charge{}, &operations.CreateCharge{
		Customer: customer,
		Amount:   amount,
		Currency: o.Currency,
	}

	err := o.client.Do(charge, payload)
	if err != nil {
		return nil, err
	}

	return charge, nil
}

func (o *Omise) DeleteCustomer(customerID string) (*omiseLib.Deletion, error) {

	model, destroy := &omiseLib.Deletion{}, &operations.DestroyCustomer{
		CustomerID: customerID,
	}

	if err := o.client.Do(model, destroy); err != nil {
		return nil, err
	}

	return model, nil
}

func (o *Omise) DeleteAllCustomer() error {

	customers, err := o.RetrieveAllCustomer()
	if err != nil {
		return err
	}

	for _, cus := range customers.Data {
		_, err = o.DeleteCustomer(cus.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
