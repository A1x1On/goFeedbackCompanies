package feedBackAux

import (
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseTripad(sParams *serviceParams, qFeedback *models.FeedbackQueryModel){
	// get All HTML
	html, err 				:= sParams.doc.Html()
	helper.CheckError(err)
	html 		  				 = strings.ToLower(html)

	// get/check title
	rateViewTripExp, err := regexp.Compile(`<span>`+ qFeedback.Company +`<\/span>.*&#34;}\" data-israteable=\"true\">`)
	helper.CheckError(err)
	listRateViewTrip     := rateViewTripExp.FindAllString(html, -1)

	for _, val := range listRateViewTrip {
		// get/set rate
		valExp, err	   := regexp.Compile(`\brating&#34;:&#34;[\d\.]*\b`)
		helper.CheckError(err)
		numExp, err	   := regexp.Compile(`\d\.?\d?$`)
		helper.CheckError(err)
		listVal  	   := valExp.FindAllString(val, -1)
		rateText 	   := numExp.FindAllString(listVal[0], -1)
		if len(rateText) != 0 {
			foldRate(rateText[0], sParams)
		}
		
		// get/set reviews
		valExp, err    = regexp.Compile(`count&#34;:&#34;\d?\b`)
		helper.CheckError(err)
		listVal        = valExp.FindAllString(val, -1)
		reviewsText   := numExp.FindAllString(listVal[0], -1)

		if len(reviewsText) != 0 {
			sParams.numReviews = getSumReviews(reviewsText[0], sParams.numReviews)
		}
	}
}
