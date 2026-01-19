package converter

import (
	models "github.com/AlexKostromin/TaskTracker/internal/tracker/domain"
	trackerV1 "github.com/AlexKostromin/TaskTracker/shared/pkg/openapi/tracker/v1"
)

func CreateTrackerRequestToModel(request *trackerV1.CreateTrackerRequest) models.CreateTrackerRequest {
	return models.CreateTrackerRequest{
		Name:        request.Name,
		Description: request.Description.Value,
	}
}
func UpdateTrackerRequestToModel(request *trackerV1.UpdateTrackerRequest) models.UpdateTrackerRequest {

	return models.UpdateTrackerRequest{
		Name:        request.Name.Or(""),
		Description: request.Description.Or(""),
	}
}
func UpdateTrackerParamsToModel(request trackerV1.UpdateTrackerParams) models.UpdateTrackerParams {

	return models.UpdateTrackerParams{
		ID: request.ID,
	}
}
