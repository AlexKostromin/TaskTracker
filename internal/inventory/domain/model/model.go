package model

import (
	"time"
)

// GetPartRequest
type GetPartRequest struct {
	UUID string
}

// GetPartResponse
type GetPartResponse struct {
	Part *Part
}

// ListPartsRequest
type ListPartsRequest struct {
	Filter *PartsFilter
}

// ListPartsResponse
type ListPartsResponse struct {
	Parts []*Part
}

// PartsFilter
type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

// Category - тип для категорий
type Category string

// Part
type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      string
	Length        float64
	Width         float64
	Height        float64
	Weight        float64
	Manufacturer  string
	Country       string
	Website       string
	Tags          []string
	Metadata      map[string]interface{}
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
