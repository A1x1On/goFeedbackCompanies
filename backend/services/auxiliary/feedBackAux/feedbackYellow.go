package feedBackAux

import (
	"gov/backend/common/helper"
	"gov/backend/models"
	"strings"
	"regexp"
)

func ParseYellow(sParams *serviceParams, qFeedback *models.FeedbackQueryModel){
	// get HTML
	html, err 	  := sParams.doc.Html()
	helper.CheckError(err)
	html 		  		= strings.ToLower(html)

	// get rate & reviews as DOM string array by found companies
	rateText      := ""
	valExp, err	  := regexp.Compile(`<div class="result-rating\s\D*\d*\)`)
	helper.CheckError(err)
	listValue  	  := valExp.FindAllString(html, -1)

	if len(listValue) != 0 {
		for _, val := range listValue {
	
			// get/set rate
			rateExp, err   := regexp.Compile(`one|two|three|four|five`)
			helper.CheckError(err)
			rateNumText    := rateExp.FindAllString(val, -1)
	
			if len(rateNumText) != 0 {
				if rateNumText[0] == "one" {
					rateText = "1.0"
				} else if rateNumText[0]  == "two" {
					rateText = "2.0"
				} else if rateNumText[0]  == "three" {
					rateText = "3.0"
				} else if rateNumText[0]  == "four" {
					rateText = "4.0"
				} else if rateNumText[0]  == "five" {
					rateText = "5.0"
				} else {
					rateText = "0.0"
				}
				foldRate(rateText, sParams)
			}
	
			// get/set reviews
			numExp, err	  := regexp.Compile(`\d*\)`)
			helper.CheckError(err)
			reviewsText   := numExp.FindAllString(val, -1)
			if len(reviewsText) != 0 {
				sParams.numReviews = getSumReviews(reviewsText[0][:1], sParams.numReviews)
			}
		}
	}

}
