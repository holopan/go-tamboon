package domain

import "time"

type Donator struct {
	ChargeID          string     `json:"charge_id"`
	CustomerID        string     `json:"customer_id"`
	Email             string     `json:"email"`
	CardID            string     `json:"card_id"`
	Number            string     `json:"number"`
	SecurityCode      string     `json:"security_code"`
	Country           string     `json:"country"`
	City              string     `json:"city"`
	Bank              string     `json:"bank"`
	PostalCode        string     `json:"postal_code"`
	Financing         string     `json:"financing"`
	LastDigits        string     `json:"last_digits"`
	Brand             string     `json:"brand"`
	ExpirationMonth   time.Month `json:"expiration_month"`
	ExpirationYear    int        `json:"expiration_year"`
	Fingerprint       string     `json:"fingerprint"`
	Name              string     `json:"name"`
	SecurityCodeCheck bool       `json:"security_code_check"`
	Amount            int64      `json:"amount"`
}

type TumboonResult struct {
	TotalDonator       int       `json:"total_donator"`
	TotalSuccessDonate int       `json:"total_success_donator"`
	Donators           []Donator `json:"donator"`
	TotalReceived      float64   `json:"total_received"`
	SuccessfulDonate   float64   `json:"successfully_donated"`
	FaultyDonation     float64   `json:"faulty_donation"`
	Average            float64   `json:"average_per_person"`
	TopDonate          []Donator `json:"top_donate"`
}
