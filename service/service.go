package service

import (
	"errors"
)

type Service struct {
}

func New() (*Service, error) {
	return nil, nil
}

func (s *Service) Run() error {
	return errors.New("Not implemented")
}
