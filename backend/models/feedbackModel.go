package models

type FeedbackModel struct {
	ServiceTitle string
	Rate         float64
	NumReviews   int
	ErrorState 	 *ErrorStateModel
}
