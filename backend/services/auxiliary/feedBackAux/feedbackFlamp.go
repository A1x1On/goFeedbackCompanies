package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
)

func ParseFlamp(sParams *serviceParams){
	errorState := &models.ErrorStateModel{Message: "null", Code: 0}
	docCount   := 0

	sParams.doc.Find("cat-brand-filial-rating").Each(func(i int, sel *goquery.Selection) {
		docCount = i + 1
		// get/set rate
		rateText, _       := sel.Attr("rating")
		foldRate(rateText, sParams)
		// get/set reviews
		reviewsText, _    := sel.Attr("reviews-count")
		sParams.numReviews = getSumReviews(reviewsText, sParams.numReviews)
	})

	// set doc html error
	if docCount == 0 {
		html, err  := sParams.doc.Html()
		helper.CheckError(err)
		setHttpErrorByHtml(html, errorState)
	}
}
