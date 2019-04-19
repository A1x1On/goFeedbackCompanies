package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseSpasibo(sParams *serviceParams, qFeedback *models.FeedbackQueryModel){
	errorState := &models.ErrorStateModel{Message: "null", Code: 0}
	docCount   := 0

	sParams.doc.Find("table.items tbody tr").Each(func(i int, sel *goquery.Selection) {
		docCount = i + 1
		// get/check title
		titleHtml      := sel.Find("td.left > .name > a")
		titleText      := strings.ToLower(titleHtml.Text())
		titleExp, err  := regexp.Compile("\\B\\s?" + qFeedback.Company  + "\"\\s\\(")
		helper.CheckError(err)

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateHtml          := sel.Find("div.stars")
			rateText, _       := rateHtml.Attr("data-fill")
			foldRate(rateText, sParams)
			// get/set reviews
			reviewsHtml       := sel.Find("td.num")
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
