package websocket

import "github.com/krobus00/websocket-service/internal/contract"

type Service struct {
	epollerService contract.EpollerService
}

func New() *Service {
	return &Service{}
}

func (s *Service) WithEpollerService(svc contract.EpollerService) *Service {
	s.epollerService = svc
	return s
}
