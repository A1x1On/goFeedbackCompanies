package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseYelp(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docCount := 0

	sParams.doc.Find(".mainAttributes__373c0__1r0QA").EachWithBreak(func(i int, sel *goquery.Selection) bool {
		docCount = i + 1
		// get/check title
		titleSel	 		:= ".heading--h3__373c0__1n4Of > a"
		titleHtml      := sel.Find(titleSel)
		// check if title tag exists
		if sel.Has(titleSel).Length() == 0 {
			setParsingErrorByCode(1000, "title" , errorState)
			return false
		}
		// ------------------------
		titleText      := strings.ToLower(titleHtml.Text())
		titleExp, err  := regexp.Compile(qFeedback.Company)
		helper.IfError(err, "can't (regexp.Compile) to get [titleExp]")

		if titleExp.MatchString(titleText) {
			// get/set rate
			rateSel	 	  := ".i-stars__373c0__30xVZ"
			rateHtml 	  := sel.Find(".i-stars__373c0__30xVZ")
			rateText, ok  := rateHtml.Attr("aria-label")
			// Tag rate can be unfound in DOM even if company has not rate. check if rate attribute exist
			if sel.Has(rateSel).Length() != 0 && !ok {
				setParsingErrorByCode(1002, "rate" , errorState)
				return false
			}
			// ------------------------
			rateExp, err  := regexp.Compile(`\d\.?\d?`)
			helper.IfError(err, "can't (rateExp.Compile) to get [rateExp]")
			listRate  	  := rateExp.FindAllString(rateText, -1)
			if len(listRate) != 0 {
				foldRate(listRate[0], sParams)
			}
		
			// get/set reviews
			reviewsSel	  := ".reviewCount__373c0__2r4xT"
			reviewsHtml   := sel.Find(reviewsSel)
			// check if reviews tag exists
			if sel.Has(rateSel).Length() != 0 && sel.Has(reviewsSel).Length() == 0 {
				setParsingErrorByCode(1000, "reviews" , errorState)
				return false
			}
			// ------------------------
			reviewsText	  := regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
			if reviewsText != "" {
				sParams.numReviews = getSumReviews(reviewsText, sParams.numReviews)
			}
			
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
