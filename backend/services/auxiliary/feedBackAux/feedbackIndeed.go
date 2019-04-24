package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseIndeed(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docFound := sParams.doc.Find(".cmp-CompanyWidget:first-child")

	// check found result by entered comapny
	if docFound.Length() == 0 {
		setParsingErrorByCode(1005, qFeedback.Company , errorState)
	}
	// ------------------------

	docFound.EachWithBreak(func(i int, sel *goquery.Selection) bool {
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
			foldRate(rateHtml.Text(), sParams, errorState)
		}

		return true
	})

}
