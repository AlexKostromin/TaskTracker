package postgres

import (
	"context"
	"database/sql"
	"fmt"

	models "github.com/AlexKostromin/TaskTracker/internal/tracker/domain"
)

type TrackerStorage struct {
	db *sql.DB
}

func NewTrackerStorage(db *sql.DB) *TrackerStorage {
	return &TrackerStorage{db: db}
}

func (s *TrackerStorage) CreateTracker(ctx context.Context, tracker models.CreateTrackerRequest) (models.Tracker, error) {
	const q = `
INSERT INTO trackers (name, description)
VALUES ($1, $2)
RETURNING id, name, COALESCE(description, '')
`
	var out models.Tracker
	if err := s.db.QueryRowContext(ctx, q, tracker.Name, nullableText(tracker.Description)).
		Scan(&out.ID, &out.Name, &out.Description); err != nil {
		return models.Tracker{}, err
	}
	return out, nil
}

func (s *TrackerStorage) UpdateTracker(ctx context.Context, request models.UpdateTrackerRequest, params models.UpdateTrackerParams) (models.Tracker, error) {
	// Пустые поля не обновляем (поведение соответствует текущему mock-репозиторию).
	const q = `
UPDATE trackers
SET
  name = COALESCE(NULLIF($1, ''), name),
  description = COALESCE(NULLIF($2, ''), description)
WHERE id = $3
RETURNING id, name, COALESCE(description, '')
`
	var out models.Tracker
	err := s.db.QueryRowContext(ctx, q, request.Name, request.Description, params.ID).
		Scan(&out.ID, &out.Name, &out.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Tracker{}, fmt.Errorf("tracker not found")
		}
		return models.Tracker{}, err
	}
	return out, nil
}

func nullableText(v string) any {
	if v == "" {
		return nil
	}
	return v
}
