package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseYell(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docFound := sParams.doc.Find("div.companies__item-content")

	// check found result by entered comapny
	if docFound.Length() == 0 {
		setParsingErrorByCode(1005, qFeedback.Company , errorState)
	}
	// ------------------------

	docFound.EachWithBreak(func(i int, sel *goquery.Selection) bool {	
		// get/check title
		titleSel	     := ".companies__item-title-text"
		titleHtml     := sel.Find(titleSel)
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
			rateSel 			   := "span.rating__value"
			rateHtml 	      := sel.Find("span.rating__value")
			// check if rate tag exists
			if sel.Has(rateSel).Length() == 0 {
				setParsingErrorByCode(1000, "rate" , errorState)
				return false
			}
			// ------------------------
			rateText 	      := trimAll(strings.ToLower(rateHtml.Text()))
			foldRate(rateText, sParams, errorState)

			// get/set reviews
			reviewsSel	      := "span.rating__reviews > span"
			reviewsHtml       := sel.Find(reviewsSel)
			// check if reviews tag exists
			if sel.Has(reviewsSel).Length() == 0 {
				setParsingErrorByCode(1000, "reviews" , errorState)
				return false
			}
			// ------------------------
			sParams.numReviews = getSumReviews(reviewsHtml.Text(), sParams.numReviews)
		}

		return true
	})

}
