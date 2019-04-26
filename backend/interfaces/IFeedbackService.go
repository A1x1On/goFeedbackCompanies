package interfaces

import (
	"gov/backend/models"
)

type IFeedbackService interface {
	GetAllByCriteria(*models.FeedbackQueryModel) []models.FeedbackModel
	GetReviewService(*models.FeedbcakParamsModel) (models.FeedbackModel, int, string)
	GetServices(int) ([]int, int, string)
	GetAll() string
}
