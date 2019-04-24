package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseSpasibo(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docFound := sParams.doc.Find("table.items tbody tr")

	// check found result by entered comapny
	if docFound.Length() == 0 {
		setParsingErrorByCode(1005, qFeedback.Company , errorState)
	}
	// ------------------------

	docFound.EachWithBreak(func(i int, sel *goquery.Selection) bool {
		// get/check title
		titleSel	 := ".left"
		titleHtml := sel.Find(titleSel)

		// check if title tag exists
		if sel.Has(titleSel).Length() == 0 {
			setParsingErrorByCode(1000, "title" , errorState)
			return false
		}
		// ------------------------
		titleText      := strings.ToLower(titleHtml.Text())
		titleExp, err  := regexp.Compile("\\B\\s?" + qFeedback.Company  + "\"\\s\\(")
		helper.IfError(err, "can't (rateExp.FindAllString) to get [titleExp]")

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateSel  := "div.stars"
			rateHtml := sel.Find(rateSel)
			// check if rate tag exists
			if sel.Has(rateSel).Length() == 0 {
				setParsingErrorByCode(1000, "rate" , errorState)
				return false
			}
			// ------------------------
			rateText, ok       := rateHtml.Attr("data-fill")
			// check if rate attribute exists
			if !ok {
				setParsingErrorByCode(1002, "rate" , errorState)
				return false
			}
			// ------------------------
			foldRate(rateText, sParams, errorState)

			// get/set reviews
			reviewsSel	      := "td.num"
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
