package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseIndeed(sParams *serviceParams, qFeedback *models.FeedbackQueryModel){
	errorState := &models.ErrorStateModel{Message: "null", Code: 0}
	docCount   := 0

	sParams.doc.Find(".cmp-CompanyWidget:first-child").Each(func(i int, sel *goquery.Selection) {
		docCount = i + 1
		// get/check title
		titleHtml     := sel.Find(".cmp-CompanyWidget-name")
		titleText     := strings.ToLower(titleHtml.Text())
		titleExp, err := regexp.Compile(qFeedback.Company)
		helper.CheckError(err)

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateHtml := sel.Find(".cmp-CompanyWidget-rating-link")
			foldRate(rateHtml.Text(), sParams)
		}
	})

	// set doc html error
	if docCount == 0 {
		html, err  := sParams.doc.Html()
		helper.CheckError(err)
		setHttpErrorByHtml(html, errorState)
	}
}
