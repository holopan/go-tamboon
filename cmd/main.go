package main

import (
	"os"

	"github.com/omise/go-tamboon/config"
	services "github.com/omise/go-tamboon/internal/core/services/csv"
	"github.com/omise/go-tamboon/internal/handlers/csv"
	"github.com/omise/go-tamboon/internal/repositories"
	"github.com/omise/go-tamboon/pkg/cipher/caesar"
)

func main() {
	omiseRepository := repositories.NewOmiseRepository(config.Setting.App.Repository.Omise.Public, config.Setting.App.Repository.Omise.Secret, config.Setting.App.Currency)

	reader := caesar.NewCaesarReader(os.Args[1], config.Setting.App.Encrypt.Caesar.Shift)
	csvService := services.NewCsvService(omiseRepository, reader, config.Setting.App.Poolsize)

	tumboonHandler := csv.NewCSVHandler(csvService)
	err := tumboonHandler.Donate()
	if err != nil {
		panic(err)
	}
}
