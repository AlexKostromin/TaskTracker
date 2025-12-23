package mock

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gitlab.com/godevs2/micro/internal/inventory/domain/model"
)

// InventoryRepositoryMock - мок-реализация репозитория для тестирования
type InventoryRepositoryMock struct {
	mu    sync.RWMutex
	parts map[string]*model.Part
}

// NewInventoryStorage создает новый мок-репозиторий с тестовыми данными
func NewInventoryStorage() *InventoryRepositoryMock {
	repo := &InventoryRepositoryMock{
		parts: make(map[string]*model.Part),
	}
	repo.initializeSampleData()
	return repo
}

// Get возвращает деталь по запросу (реализация InventoryStorage.Get)
func (r *InventoryRepositoryMock) Get(_ context.Context, req *model.GetPartRequest) (*model.GetPartResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	part := r.parts[req.UUID]

	return &model.GetPartResponse{
		Part: part,
	}, nil
}

// ListParts возвращает список деталей с фильтрацией (реализация InventoryStorage.ListParts)
func (r *InventoryRepositoryMock) ListParts(_ context.Context, req *model.ListPartsRequest) (*model.ListPartsResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*model.Part

	// Если фильтр пустой, возвращаем все детали
	if req == nil || req.Filter == nil {
		for _, part := range r.parts {
			result = append(result, part)
		}
		return &model.ListPartsResponse{
			Parts: result,
		}, nil
	}

	// Фильтруем детали
	for _, part := range r.parts {
		if r.matchesFilter(part, req.Filter) {
			result = append(result, part)
		}
	}

	return &model.ListPartsResponse{
		Parts: result,
	}, nil
}

// matchesFilter проверяет, соответствует ли деталь фильтру
func (r *InventoryRepositoryMock) matchesFilter(part *model.Part, filter *model.PartsFilter) bool {
	if filter == nil {
		return true
	}

	// Фильтр по UUID
	if len(filter.UUIDs) > 0 {
		found := false
		for _, uuid := range filter.UUIDs {
			if part.UUID == uuid {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Фильтр по названиям
	if len(filter.Names) > 0 {
		found := false
		for _, name := range filter.Names {
			if part.Name == name {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Фильтр по категориям
	if len(filter.Categories) > 0 {
		found := false
		for _, category := range filter.Categories {
			if string(category) == part.Category {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Фильтр по странам производителей
	if len(filter.ManufacturerCountries) > 0 {
		found := false
		for _, country := range filter.ManufacturerCountries {
			if part.Country == country {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Фильтр по тегам
	if len(filter.Tags) > 0 {
		tagFound := false
		for _, filterTag := range filter.Tags {
			for _, partTag := range part.Tags {
				if partTag == filterTag {
					tagFound = true
					break
				}
			}
			if tagFound {
				break
			}
		}
		if !tagFound {
			return false
		}
	}

	return true
}

// initializeSampleData инициализирует тестовые данные
func (r *InventoryRepositoryMock) initializeSampleData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Двигатель
	r.parts["engine_001"] = &model.Part{
		UUID:          "engine_001",
		Name:          "Ионный двигатель Mark IV",
		Description:   "Высокоэффективный ионный двигатель для межпланетных перелетов",
		Price:         2500000.99,
		StockQuantity: 5,
		Category:      "CATEGORY_ENGINE",
		Length:        3.5,
		Width:         2.1,
		Height:        2.1,
		Weight:        4500.0,
		Manufacturer:  "Quantum Propulsion Systems",
		Country:       "USA",
		Website:       "https://quantumprop.com",
		Tags:          []string{"engine", "propulsion", "ion", "high-efficiency"},
		Metadata: map[string]interface{}{
			"thrust":         150000,
			"fuel_type":      "xenon",
			"lifespan_years": 20,
			"warranty":       true,
		},
		CreatedAt: now.Add(-30 * 24 * time.Hour),
		UpdatedAt: now,
	}

	// Иллюминатор
	r.parts["porthole_002"] = &model.Part{
		UUID:          "porthole_002",
		Name:          "Титановый иллюминатор",
		Description:   "Иллюминатор с бронированным стеклом для космических кораблей",
		Price:         125000.50,
		StockQuantity: 15,
		Category:      "CATEGORY_PORTHOLE",
		Length:        1.2,
		Width:         1.2,
		Height:        0.3,
		Weight:        85.5,
		Manufacturer:  "Stellar Optics",
		Country:       "Germany",
		Website:       "https://stellar-optics.de",
		Tags:          []string{"window", "porthole", "armored", "viewing"},
		Metadata: map[string]interface{}{
			"diameter_cm":     120,
			"glass_type":      "armored_crystal",
			"pressure_rating": "10atm",
			"uv_protection":   true,
		},
		CreatedAt: now.Add(-15 * 24 * time.Hour),
		UpdatedAt: now,
	}

	// Топливный бак
	r.parts["fuel_tank_003"] = &model.Part{
		UUID:          "fuel_tank_003",
		Name:          "Криогенный топливный бак",
		Description:   "Бак для хранения жидкого водорода и кислорода",
		Price:         850000.00,
		StockQuantity: 3,
		Category:      "CATEGORY_FUEL",
		Length:        4.8,
		Width:         2.4,
		Height:        2.4,
		Weight:        3200.0,
		Manufacturer:  "CryoTech Solutions",
		Country:       "Japan",
		Website:       "https://cryotech.jp",
		Tags:          []string{"fuel", "tank", "cryogenic", "storage"},
		Metadata: map[string]interface{}{
			"capacity_l":   5000,
			"material":     "titanium_composite",
			"insulation":   "vacuum_multilayer",
			"max_pressure": 50,
		},
		CreatedAt: now.Add(-45 * 24 * time.Hour),
		UpdatedAt: now.Add(-7 * 24 * time.Hour),
	}

	// Крыло
	r.parts["wing_004"] = &model.Part{
		UUID:          "wing_004",
		Name:          "Композитное крыло",
		Description:   "Крыло из углепластика для атмосферного полета",
		Price:         1750000.00,
		StockQuantity: 2,
		Category:      "CATEGORY_WING",
		Length:        12.5,
		Width:         4.2,
		Height:        1.8,
		Weight:        2800.0,
		Manufacturer:  "AeroComposite Ltd",
		Country:       "UK",
		Website:       "https://aerocomposite.co.uk",
		Tags:          []string{"wing", "composite", "aerodynamic", "carbon_fiber"},
		Metadata: map[string]interface{}{
			"span_m":      25.0,
			"area_sqm":    52.5,
			"material":    "carbon_fiber_epoxy",
			"max_load_kg": 8500,
		},
		CreatedAt: now.Add(-60 * 24 * time.Hour),
		UpdatedAt: now.Add(-14 * 24 * time.Hour),
	}

	log.Printf("📊 Инициализировано %d тестовых деталей в репозитории", len(r.parts))
}
