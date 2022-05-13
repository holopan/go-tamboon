package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type (
	Conf struct {
		App struct {
			Encrypt    Encrypt    `mapstructure:"encrypt"`
			Repository Repository `mapstructure:"repository"`
			Currency   string     `mapstructure:"currency"`
			Poolsize   int        `mapstructure:"pool_size"`
		} `mapstructure:"app"`
	}

	Encrypt struct {
		Caesar Caesar `mapstructure:"caesar"`
	}

	Caesar struct {
		Shift int `mapstructure:"shift"`
	}

	Repository struct {
		Omise Omise `mapstructure:"omise"`
	}

	Omise struct {
		Public string `mapstructure:"public"`
		Secret string `mapstructure:"secret"`
	}
)

// Setting is instance config setting
var Setting Conf

// Stage is environment program
var Stage string

func init() {
	setup()
	log.Printf(" [*] Env name : %s\n", Stage)
}

func setup() {
	Stage = os.Getenv("ENVIRONMENT")
	if Stage == "" {
		Stage = "development"
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(fmt.Sprintf(".env.%s", Stage))
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Setting)
	if err != nil {
		panic(err)
	}
}
