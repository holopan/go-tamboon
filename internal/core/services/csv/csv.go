package services

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/omise/go-tamboon/internal/core/domain"
	"github.com/omise/go-tamboon/internal/core/ports"
	"github.com/omise/go-tamboon/pkg/cipher/caesar"
)

type service struct {
	mu               sync.Mutex
	donateRepository ports.DonateRepository
	Reader           *caesar.CaesarReader
	PoolSize         int
}

type data struct {
	Line   string
	Result *domain.TumboonResult
}

func NewCsvService(donateRepository ports.DonateRepository, reader *caesar.CaesarReader, poolSize int) *service {
	return &service{
		donateRepository: donateRepository,
		Reader:           reader,
		PoolSize:         poolSize,
	}
}

func (sv *service) Donate() (*domain.TumboonResult, error) {
	var count int
	result := domain.TumboonResult{
		TopDonate: make([]domain.Donator, 3),
	}

	c := make(chan data, 5)

	wg := sync.WaitGroup{}

	for w := 1; w <= 5; w++ {
		go sv.process(w, c, &wg)
	}

	for {
		count++
		line, ok := sv.Reader.ReadLine()

		if ok == false {
			break
		}
		if count == 1 {
			continue
		}
		input := data{
			Line:   line,
			Result: &result,
		}
		wg.Add(1)
		c <- input

	}
	time.Sleep(5 * time.Second)
	close(c)

	wg.Wait()
	return &result, nil
}

func (sv *service) process(id int, c chan data, wg *sync.WaitGroup) {
	for input := range c {
		var successFullyDonate, faultyDonation, amount int64
		var expireYear, expireMonth, success int
		var err error
		var result *domain.TumboonResult
		errFlag := false

		result = input.Result

		row := strings.Split(input.Line, ",")

		amount, err = strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			// log.Printf("Donator %s cannot convert amount %s to int", row[0], row[1])
			errFlag = true
		}

		expireMonth, err = strconv.Atoi(row[4])
		if err != nil {
			// log.Printf("Donator %s cannot convert exp month %s to int", row[0], row[4])
			faultyDonation = amount
			errFlag = true
		}

		expireYear, err = strconv.Atoi(row[5])
		if err != nil {
			// log.Printf("Donator %s cannot convert exp year %s to int", row[0], row[5])
			faultyDonation = amount
			errFlag = true
		}
		if expireYear < time.Now().Year() {
			// log.Printf("Donator %s exp year %s is in the past", row[0], row[5])
			faultyDonation = amount
			errFlag = true

		} else if expireYear == time.Now().Year() && expireMonth <= int(time.Now().Month()) {
			// log.Printf("Donator %s exp month %d and year %d is in the past", row[0], expireMonth, expireYear)
			faultyDonation = amount
			errFlag = true
		}

		donator := domain.Donator{
			Name:            row[0],
			Amount:          amount,
			Number:          row[2],
			SecurityCode:    row[3],
			ExpirationMonth: time.Month(expireMonth),
			ExpirationYear:  expireYear,
		}

		if !errFlag {
			err = sv.donateRepository.Create(&donator)
			if err != nil {
				errFlag = false
				faultyDonation = donator.Amount
				// log.Printf("Donator %s cannot Create record because %s", donator.Name, err.Error())
			} else {
				successFullyDonate = amount
				success = 1
			}
		}

		result.Donators = append(result.Donators, donator)

		// protect Race condition
		sv.mu.Lock()
		result.SuccessfulDonate = result.SuccessfulDonate + (float64(successFullyDonate) / 100)
		result.FaultyDonation = result.FaultyDonation + (float64(faultyDonation) / 100)
		result.TotalReceived = result.TotalReceived + (float64(amount) / 100)
		result.TotalDonator = result.TotalDonator + 1
		result.TotalSuccessDonate = result.TotalSuccessDonate + success

		if !errFlag {
			sv.addTopDonate(donator, result)
		}

		sv.mu.Unlock()
		wg.Done()
	}

}

func (sv *service) addTopDonate(newDonator domain.Donator, result *domain.TumboonResult) {
	donator := newDonator
	for i, d := range result.TopDonate {
		if d.Name == "" {
			result.TopDonate[i] = donator
			break
		} else if d.Amount <= donator.Amount {
			result.TopDonate[i] = donator
			donator = d
			continue
		} else {
			continue
		}
	}
}
