package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParsePravda(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docCount := 0

	sParams.doc.Find(".mdc-companies-item-title").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		docCount   = i + 1

		// get/check title
		titleSel	 := "span > a"
		titleHtml := sel.Find("span > a")
		// check if title tag exists
		if sel.Has(titleSel).Length() == 0 {
			setParsingErrorByCode(1000, "title" , errorState)
			return false
		}
		// ------------------------
		titleText     := strings.ToLower(titleHtml.Text())
		titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
		helper.IfError(err, "can't (rateExp.Compile) to get [titleExp]")

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateSel  := ".mdc-companies-item-rating > span"
			rateHtml := sel.Find(rateSel)
			// check if rate tag exists
			if sel.Has(rateSel).Length() == 0 {
				setParsingErrorByCode(1000, "rate" , errorState)
				return false
			}
			// ------------------------
			rateText, ok := rateHtml.Attr("data-rating")
			// check if rate attribute exists
			if !ok {
				setParsingErrorByCode(1002, "rate" , errorState)
				return false
			}
			// ------------------------
			foldRate(rateText, sParams)
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
