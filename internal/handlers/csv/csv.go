package csv

import (
	"log"

	"github.com/omise/go-tamboon/internal/core/ports"
)

type CSVHandler struct {
	tumboonService ports.TumboonService
}

func NewCSVHandler(tumboonService ports.TumboonService) *CSVHandler {
	return &CSVHandler{
		tumboonService: tumboonService,
	}
}

func (h *CSVHandler) Donate() error {
	log.Println("performing donations...")

	result, err := h.tumboonService.Donate()
	if err != nil {
		return err
	}

	log.Println("done.")
	log.Printf("total donator: %d\n", result.TotalDonator)
	log.Printf("total Success donate: %d\n", result.TotalSuccessDonate)
	log.Println("")
	log.Printf("	 total received: THB  %.2f\n", result.TotalReceived)
	log.Printf("	 successfully donated: THB  %.2f\n", result.SuccessfulDonate)
	log.Printf("	 faulty donation: THB  %.2f\n", result.FaultyDonation)
	log.Printf("	 average per person: THB  %.2f\n", result.SuccessfulDonate/float64(result.TotalSuccessDonate))
	log.Println("		  top donators:")
	for _, donator := range result.TopDonate {
		amount := float64(donator.Amount)
		amount = amount / 100
		log.Printf("		  Name : %s Amount : %.2f", donator.Name, amount)
	}

	return nil
}
