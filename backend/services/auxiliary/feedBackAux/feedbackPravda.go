package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParsePravda(sParams *serviceParams, qFeedback *models.FeedbackQueryModel){
	errorState := &models.ErrorStateModel{Message: "null", Code: 0}
	docCount   := 0

	sParams.doc.Find("div.w_star").Each(func(i int, sel *goquery.Selection) {
		docCount = i + 1
		// get/check title
		titleHtml      := sel.Find(".m_title > .flw a")
		titleText      := strings.ToLower(titleHtml.Text())
		titleExp, err  := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
		helper.CheckError(err)

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateHtml          := sel.Find(".img_p > .p_f_s")
			rateExp, err      := regexp.Compile("[\\d\\.]*")
			rateText          := rateExp.FindAllString(trimAll(rateHtml.Text()), -1)
			helper.CheckError(err)
			foldRate(rateText[0], sParams)
			// get/set reviews
			reviewsHtml       := sel.Find(".img_p > .numReviews")
			reviewsText	      := regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
			sParams.numReviews = getSumReviews(reviewsText, sParams.numReviews)
		}
	})

	// set doc html error
	if docCount == 0 {
		html, err  := sParams.doc.Html()
		helper.CheckError(err)
		setHttpErrorByHtml(html, errorState)
	}
}
