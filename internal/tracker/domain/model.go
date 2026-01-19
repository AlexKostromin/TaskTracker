package models

type Tracker struct {
	ID          int
	Name        string
	Description string
}

type CreateTrackerRequest struct {
	Name        string
	Description string
}

type UpdateTrackerParams struct {
	ID int
}

type UpdateTrackerRequest struct {
	Name        string
	Description string
}
