package models

type FeedbackQueryModel struct {
	Country      string
	Company      string 
	ISOCode	    string
	AvarageRate  float64
	NumReviews   int
	Services     []*FeedbackServiceModel
}
