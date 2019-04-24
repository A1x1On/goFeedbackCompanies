package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseApoi(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docFound := sParams.doc.Find("div.w_star")

	// check found result by entered comapny
	if docFound.Length() == 0 {
		setParsingErrorByCode(1005, qFeedback.Company , errorState)
	}
	// ------------------------

	docFound.EachWithBreak(func(i int, sel *goquery.Selection) bool {
		// get/check title
		titleSel	 := ".m_title > .flw a"
		titleHtml := sel.Find(titleSel)
		// check if title tag exists
		if sel.Has(titleSel).Length() == 0 {
			setParsingErrorByCode(1000, "title" , errorState)
			return false
		}
		// ------------------------
		titleText      := strings.ToLower(titleHtml.Text())
		titleExp, err  := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
		helper.IfError(err, "can't (regexp.Compile) to get [titleExp]")

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateSel	 			:= ".img_p > .p_f_s"
			rateHtml  			:= sel.Find(rateSel)
			// check if rate tag exists
			if sel.Has(rateSel).Length() == 0 {
				setParsingErrorByCode(1000, "rate" , errorState)
				return false
			}
			// ------------------------
			rateExp, err      := regexp.Compile("[\\d\\.]*")
			rateText          := rateExp.FindAllString(trimAll(rateHtml.Text()), -1)
			helper.IfError(err, "can't (rateExp.FindAllString) to get [rateText]")
			foldRate(rateText[0], sParams, errorState)

			// get/set reviews
			reviewsSel	      := ".img_p > .numReviews"
			reviewsHtml       := sel.Find(reviewsSel)
			// check if reviews tag exists
			if sel.Has(reviewsSel).Length() == 0 {
				setParsingErrorByCode(1000, "reviews" , errorState)
				return false
			}
			// ------------------------
			reviewsText	  		:= regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
			sParams.numReviews = getSumReviews(reviewsText, sParams.numReviews)
		}

		return true
	})
		
}
