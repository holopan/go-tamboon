package ports

import "github.com/omise/go-tamboon/internal/core/domain"

type TumboonService interface {
	Donate() (*domain.TumboonResult, error)
}
