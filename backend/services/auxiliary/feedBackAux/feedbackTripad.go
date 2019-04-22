package feedBackAux

import (
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseTripad(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	docCount				   := 0

	// get All HTML
	html, err 				:= sParams.doc.Html()
	helper.IfError(err, "can't (sParams.doc.Html) to get [html]")
	html 		  				 = strings.ToLower(html)

	// get/check title
	rateViewTripExp, err := regexp.Compile(`<span>`+ qFeedback.Company +`<\/span>.*&#34;}\" data-israteable=\"true\">`)
	helper.IfError(err, "can't (regexp.Compile) to get [rateViewTripExp]")
	listRateViewTrip     := rateViewTripExp.FindAllString(html, -1)

	for i, val := range listRateViewTrip {
		docCount = i + 1

		// get/set rate
		valExp, err	   := regexp.Compile(`\brating&#34;:&#34;[\d\.]*\b`)
		helper.IfError(err, "can't (regexp.Compile) to get [valExp]")
		numExp, err	   := regexp.Compile(`\d\.?\d?$`)
		helper.IfError(err, "can't (regexp.Compile) to get [numExp]")
		listVal  	   := valExp.FindAllString(val, -1)
		rateText 	   := numExp.FindAllString(listVal[0], -1)
		if len(rateText) != 0 {
			foldRate(rateText[0], sParams)
		}
		
		// get/set reviews
		valExp, err    = regexp.Compile(`count&#34;:&#34;\d?\b`)
		helper.IfError(err, "can't (regexp.Compile) to get [valExp]")
		listVal        = valExp.FindAllString(val, -1)
		reviewsText   := numExp.FindAllString(listVal[0], -1)

		if len(reviewsText) != 0 {
			sParams.numReviews = getSumReviews(reviewsText[0], sParams.numReviews)
		}
	}

	// set doc html error
	if docCount == 0 {
		html, err  := sParams.doc.Html()
		helper.IfError(err, "can't (sParams.doc.Html()) to get [html]")
		setHttpErrorByHtml(html, errorState)
	}
}
