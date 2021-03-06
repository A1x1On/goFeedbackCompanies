package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseOtzyv(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docFound := sParams.doc.Find(".otzyv_box_float")

	// check found result by entered comapny
	if docFound.Length() == 0 {
		setParsingErrorByCode(1005, qFeedback.Company , errorState)
	}
	// ------------------------

	docFound.EachWithBreak(func(i int, sel *goquery.Selection) bool {
		valSel   := ".otzyv_item_cat1"
		valHtml  := sel.Find(valSel)
		// check if main tag exists
		if sel.Has(valSel).Length() == 0 {
			setParsingErrorByCode(1000, "main block" , errorState)
			return false
		}
		// ------------------------

		// get/set rate
		rateExp, err   := regexp.Compile(`г\s\d\.?\d?`)
		helper.IfError(err, "can't (regexp.Compile) to get [rateExp]")
		listRate 	   := rateExp.FindAllString(valHtml.Text(), -1)

		if len(listRate) != 0 {
			rateExp, err  := regexp.Compile(`\d\.?\d?`)
			helper.IfError(err, "can't (regexp.Compile) to get [rateExp]")
			listRate  	  := rateExp.FindAllString(trimAll(listRate[0]), -1)
			if len(listRate) != 0 {
				foldRate(listRate[0], sParams, errorState)
			}
		}

		// get/set reviews
		reviewsExp, err := regexp.Compile(`\s\d*\sо`)
		helper.IfError(err, "can't (regexp.Compile) to get [reviewsExp]")
		listReviews     := reviewsExp.FindAllString(valHtml.Text(), -1)
		if len(listReviews) != 0 {
			reviewsText			:= regexp.MustCompile("[\\D]*").ReplaceAllString(strings.TrimLeft(listReviews[0], " "), "")
			sParams.numReviews = getSumReviews(reviewsText, sParams.numReviews)
		}

		return true
	})

}
