package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseYell(sParams *serviceParams, qFeedback *models.FeedbackQueryModel){
	errorState := &models.ErrorStateModel{Message: "null", Code: 0}
	docCount   := 0

	sParams.doc.Find("div.companies__item-content").Each(func(i int, sel *goquery.Selection) {	
		docCount = i + 1
		// get/check title
		titleHtml     := sel.Find(".companies__item-title-text")
		titleText     := strings.ToLower(titleHtml.Text())
		titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
		helper.CheckError(err)

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateHtml 	      := sel.Find("span.rating__value")
			rateText 	      := trimAll(strings.ToLower(rateHtml.Text()))
			foldRate(rateText, sParams)
			// get/set reviews
			reviewsHtml       := sel.Find("span.rating__reviews > span")
			sParams.numReviews = getSumReviews(reviewsHtml.Text(), sParams.numReviews)
		}
	})

	// set doc html error
	if docCount == 0 {
		html, err  := sParams.doc.Html()
		helper.CheckError(err)
		setHttpErrorByHtml(html, errorState)
	}
}
