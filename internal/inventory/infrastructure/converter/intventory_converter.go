package converter

import (
	"fmt"

	"gitlab.com/godevs2/micro/internal/inventory/domain/model"
	inventoryV1 "gitlab.com/godevs2/micro/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetPartRequestToModel конвертирует gRPC запрос в доменную модель с валидацией
func GetPartRequestToModel(req *inventoryV1.GetPartRequest) (*model.GetPartRequest, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.Uuid == "" {
		return nil, fmt.Errorf("uuid is required")
	}
	return &model.GetPartRequest{
		UUID: req.Uuid,
	}, nil
}

// GetPartResponseToProto конвертирует доменный ответ в gRPC ответ
func GetPartResponseToProto(resp *model.GetPartResponse) (*inventoryV1.GetPartResponse, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	return &inventoryV1.GetPartResponse{
		Part: PartToProto(resp.Part),
	}, nil
}

// ListPartsRequestToModel конвертирует gRPC запрос в доменную модель с обработкой ошибок
func ListPartsRequestToModel(req *inventoryV1.ListPartsRequest) (*model.ListPartsRequest, *model.ListPartsRequest) {
	if req == nil {
		return &model.ListPartsRequest{}, nil
	}

	filter, err := PartsFilterToModel(req.GetFilter())
	if err != nil {
		return nil, nil
	}

	return &model.ListPartsRequest{
		Filter: filter,
	}, nil
}

// PartsFilterToModel конвертирует gRPC фильтр в доменный фильтр с обработкой категорий
func PartsFilterToModel(filter *inventoryV1.PartsFilter) (*model.PartsFilter, error) {
	if filter == nil {
		return nil, nil
	}

	// Конвертируем категории с обработкой ошибок
	var categories []model.Category
	for _, cat := range filter.GetCategories() {
		// Преобразуем proto enum в строку (название enum)
		// Например, CATEGORY_ENGINE -> "CATEGORY_ENGINE"
		categories = append(categories, model.Category(cat.String()))
	}

	return &model.PartsFilter{
		UUIDs:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            categories,
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}, nil
}

// ListPartsResponseToProto конвертирует доменный ответ в gRPC ответ
func ListPartsResponseToProto(resp *model.ListPartsResponse) *inventoryV1.ListPartsResponse {
	if resp == nil {
		return &inventoryV1.ListPartsResponse{}
	}

	// Конвертируем каждый Part из доменной модели в proto
	var protoParts []*inventoryV1.Part
	for _, part := range resp.Parts {
		protoParts = append(protoParts, PartToProto(part))
	}

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}
}

// PartToProto конвертирует доменную Part в gRPC Part
func PartToProto(part *model.Part) *inventoryV1.Part {
	if part == nil {
		return nil
	}

	var createdAt, updatedAt *timestamppb.Timestamp
	if !part.CreatedAt.IsZero() {
		createdAt = timestamppb.New(part.CreatedAt)
	}
	if !part.UpdatedAt.IsZero() {
		updatedAt = timestamppb.New(part.UpdatedAt)
	}

	// Конвертируем категорию из строки в proto enum
	var category inventoryV1.Category
	switch part.Category {
	case "CATEGORY_ENGINE":
		category = inventoryV1.Category_CATEGORY_ENGINE
	case "CATEGORY_FUEL":
		category = inventoryV1.Category_CATEGORY_FUEL
	case "CATEGORY_PORTHOLE":
		category = inventoryV1.Category_CATEGORY_PORTHOLE
	case "CATEGORY_WING":
		category = inventoryV1.Category_CATEGORY_WING
	default:
		category = inventoryV1.Category_CATEGORY_UNSPECIFIED
	}

	// Конвертируем метаданные
	protoMetadata := MetadataToProto(part.Metadata)

	return &inventoryV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      category,
		Dimensions: &inventoryV1.Dimensions{
			Length: part.Length,
			Width:  part.Width,
			Height: part.Height,
			Weight: part.Weight,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    part.Manufacturer,
			Country: part.Country,
			Website: part.Website,
		},
		Tags:      part.Tags,
		Metadata:  protoMetadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// MetadataToProto конвертирует метаданные из доменного формата в proto
func MetadataToProto(metadata map[string]interface{}) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	protoMetadata := make(map[string]*inventoryV1.Value)
	for k, v := range metadata {
		val := &inventoryV1.Value{}
		switch v := v.(type) {
		case string:
			val.Value = &inventoryV1.Value_StringValue{StringValue: v}
		case int:
			val.Value = &inventoryV1.Value_Int64Value{Int64Value: int64(v)}
		case int64:
			val.Value = &inventoryV1.Value_Int64Value{Int64Value: v}
		case float32:
			val.Value = &inventoryV1.Value_DoubleValue{DoubleValue: float64(v)}
		case float64:
			val.Value = &inventoryV1.Value_DoubleValue{DoubleValue: v}
		case bool:
			val.Value = &inventoryV1.Value_BoolValue{BoolValue: v}
		default:
			// Если тип не поддерживается, пропускаем или преобразуем в строку
			continue
		}
		protoMetadata[k] = val
	}
	return protoMetadata
}
