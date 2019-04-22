package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseIndeed(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docCount   := 0

	sParams.doc.Find(".cmp-CompanyWidget:first-child").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		docCount 		= i + 1
		// get/check title
		titleSel	 	  := ".cmp-CompanyWidget-name"
		titleHtml     := sel.Find(titleSel)
		// check if title tag exists
		if sel.Has(titleSel).Length() == 0 {
			setParsingErrorByCode(1000, "title" , errorState)
			return false
		}
		// ------------------------
		titleText     := strings.ToLower(titleHtml.Text())
		titleExp, err := regexp.Compile(qFeedback.Company)
		helper.IfError(err, "can't (regexp.Compile) to get [titleExp]")

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateSel	:= ".cmp-CompanyWidget-rating-link"
			rateHtml := sel.Find(rateSel)
			// check if rate tag exists
			if sel.Has(rateSel).Length() == 0 {
				setParsingErrorByCode(1000, "rate" , errorState)
				return false
			}
			// ------------------------
			foldRate(rateHtml.Text(), sParams)
		}

		return true
	})

	// set doc html error
	if docCount == 0 {
		html, err  := sParams.doc.Html()
		helper.IfError(err, "can't (sParams.doc.Html()) to get [html]")
		setHttpErrorByHtml(html, errorState)
	}
}
