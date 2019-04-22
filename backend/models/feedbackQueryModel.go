package models

type FeedbackQueryModel struct {
	Company      string 
	ISOCode	    string
	AvarageRate  float64
	NumReviews   int
	Services     []*FeedbackServiceModel
}
