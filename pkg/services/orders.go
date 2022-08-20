package services

import (
	"context"

	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/pkg/errors"
)

type Orders struct{}

func NewOrders() *Orders {
	return &Orders{}
}

func (s *Orders) Create(ctx context.Context, o types.Order) (string, error) {
	return "", errors.New("not implemented")
}

func (s *Orders) Get(ctx context.Context, id string) (types.Order, error) {
	return types.Order{}, errors.New("not implemented")
}

func (s *Orders) List(ctx context.Context, limit, offset int) ([]types.Order, error) {
	return nil, errors.New("not implemented")
}

func (s *Orders) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
