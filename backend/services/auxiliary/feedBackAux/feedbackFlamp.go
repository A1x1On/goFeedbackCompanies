package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
)

func ParseFlamp(sParams *serviceParams, errorState *models.ErrorStateModel){
	docCount := 0

	sParams.doc.Find("cat-brand-filial-rating").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		docCount 	  = i + 1

		// get/set rate
		rateText, ok := sel.Attr("rating")
		// check if rate tag exists
		if !ok {
			setParsingErrorByCode(1002, "rate" , errorState)
			return false
		}
		// ------------------------
		foldRate(rateText, sParams)

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

	// set doc html error
	if docCount == 0 {
		html, err  := sParams.doc.Html()
		helper.IfError(err, "can't (sParams.doc.Html()) to get [html]")
		setHttpErrorByHtml(html, errorState)
	}
}
