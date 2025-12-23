package grpc_server

import (
	"context"

	"gitlab.com/godevs2/micro/internal/inventory/infrastructure/converter"
	inventoryV1 "gitlab.com/godevs2/micro/shared/pkg/proto/inventory/v1"
)

func (s *Server) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	// Конвертируем запрос
	domainReq, err := converter.GetPartRequestToModel(req)
	if err != nil {
		return nil, err
	}

	// Вызываем сервис - Get возвращает 2 значения
	domainResp, err := s.inventoryService.Get(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Конвертируем ответ
	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(domainResp.Part),
	}, nil
}
