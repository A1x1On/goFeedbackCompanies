package interfaces

import (
	"gov/backend/models"
)

type IFeedbackService interface {
	SetRates(*models.FeedbackQueryModel, *[]models.FeedbackModel)
}
