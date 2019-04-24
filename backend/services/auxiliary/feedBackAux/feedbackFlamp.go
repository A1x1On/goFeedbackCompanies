package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/models"
)

func ParseFlamp(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docFound := sParams.doc.Find("cat-brand-filial-rating")

	// check found result by entered comapny
	if docFound.Length() == 0 {
		setParsingErrorByCode(1005, qFeedback.Company , errorState)
	}
	// ------------------------

	docFound.EachWithBreak(func(i int, sel *goquery.Selection) bool {
		// get/set rate
		rateText, ok := sel.Attr("rating")
		// check if rate tag exists
		if !ok {
			setParsingErrorByCode(1002, "rate" , errorState)
			return false
		}
		// ------------------------
		foldRate(rateText, sParams, errorState)

		// get/set reviews
		reviewsText, ok    := sel.Attr("reviews-count")
		// check if reviews attribute exists
		if !ok {
			setParsingErrorByCode(1002, "reviews" , errorState)
			return false
		}
		// ------------------------
		sParams.numReviews = getSumReviews(reviewsText, sParams.numReviews)

		return true
	})

}
