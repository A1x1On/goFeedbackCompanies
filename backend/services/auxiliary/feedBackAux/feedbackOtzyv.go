package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseOtzyv(sParams *serviceParams, qFeedback *models.FeedbackQueryModel){
	errorState := &models.ErrorStateModel{Message: "null", Code: 0}
	docCount   := 0

	sParams.doc.Find(".otzyv_box_float").Each(func(i int, sel *goquery.Selection) {
		docCount = i + 1
		valHtml        := sel.Find(".otzyv_item_cat1")

		// get/set rate
		rateExp, err   := regexp.Compile(`г\s\d\.?\d?`)
		helper.CheckError(err)
		listRate 	   := rateExp.FindAllString(valHtml.Text(), -1)

		if len(listRate) != 0 {
			rateExp, err  := regexp.Compile(`\d\.?\d?`)
			helper.CheckError(err)
			listRate  	  := rateExp.FindAllString(trimAll(listRate[0]), -1)
			if len(listRate) != 0 {
				foldRate(listRate[0], sParams)
			}
		}

		// get/set reviews
		reviewsExp, err := regexp.Compile(`\s\d*\sо`)
		helper.CheckError(err)
		listReviews     := reviewsExp.FindAllString(valHtml.Text(), -1)
		if len(listReviews) != 0 {
			reviewsText			:= regexp.MustCompile("[\\D]*").ReplaceAllString(strings.TrimLeft(listReviews[0], " "), "")
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
