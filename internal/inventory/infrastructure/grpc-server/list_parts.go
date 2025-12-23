package grpc_server

import (
	"context"

	"gitlab.com/godevs2/micro/internal/inventory/infrastructure/converter"
	inventoryV1 "gitlab.com/godevs2/micro/shared/pkg/proto/inventory/v1"
)

func (s *Server) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	// Конвертируем gRPC запрос в доменную модель
	domainReq, _ := converter.ListPartsRequestToModel(req)

	// Вызываем бизнес-логику (application сервис)
	// Обратите внимание: ListParts возвращает 2 значения - response и error
	domainResp, err := s.inventoryService.ListParts(ctx, domainReq)
	if err != nil {
		return nil, err
	}

	// Конвертируем доменный ответ обратно в gRPC ответ
	return converter.ListPartsResponseToProto(domainResp), nil
}
