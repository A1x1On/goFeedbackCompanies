package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseYelp(sParams *serviceParams, qFeedback *models.FeedbackQueryModel){
	errorState := &models.ErrorStateModel{Message: "null", Code: 0}
	docCount   := 0

	sParams.doc.Find(".mainAttributes__373c0__1r0QA").Each(func(i int, sel *goquery.Selection) {
		docCount = i + 1
		// get/check title
		titleHtml      := sel.Find(".heading--h3__373c0__1n4Of > a")
		titleText      := strings.ToLower(titleHtml.Text())
		titleExp, err  := regexp.Compile(qFeedback.Company)
		helper.CheckError(err)

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateHtml 	  := sel.Find(".i-stars__373c0__30xVZ")
			rateText, _   := rateHtml.Attr("aria-label")
			rateExp, err  := regexp.Compile(`\d\.?\d?`)
			helper.CheckError(err)
			listRate  	  := rateExp.FindAllString(rateText, -1)
			if len(listRate) != 0 {
				foldRate(listRate[0], sParams)
			}
		
			// get/set reviews
			reviewsHtml   := sel.Find(".reviewCount__373c0__2r4xT")
			reviewsText	  := regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
			if reviewsText != "" {
				sParams.numReviews = getSumReviews(reviewsText, sParams.numReviews)
			}
			
		}
	})

	// set doc html error
	if docCount == 0 {
		html, err  := sParams.doc.Html()
		helper.CheckError(err)
		setHttpErrorByHtml(html, errorState)
	}
}
