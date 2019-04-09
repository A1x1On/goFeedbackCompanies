package interfaces

import (
	"gov/backend/models"
)

type IFeedbackService interface {
	GetAllByCriteria(*models.FeedbackQueryModel) []models.FeedbackModel
}
