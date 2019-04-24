package feedBackAux

import (
	"gov/backend/common/helper"
	"gov/backend/models"
	"encoding/json"
	"strings"
	"regexp"
)

type bbbRating struct {
	Search struct {
		Results []struct {
			BusinessName        string `json:"businessName"`
			Rating              string `json:"rating"`
		} `json:"results"`
	} `json:"search"`
}

func ParseBBB(sParams *serviceParams, errorState *models.ErrorStateModel, qFeedback *models.FeedbackQueryModel){
	// get HTML
	html, err 	  := sParams.doc.Html()
	helper.IfError(err, "can't (sParams.doc.Html) to get [html]")
	html 		  		= strings.ToLower(html)

	// get title & rate as json by found companies
	sRating  	  := bbbRating{}
	rateText      := ""
	valExp, err	  := regexp.Compile(`bbbDtmData.*\}`)
	helper.IfError(err, "can't (regexp.Compile) to get [valExp]")
	listValue  	  := valExp.FindAllString(sParams.doc.Text(), -1)

	if len(listValue) != 0 {
		jsonExp, err := regexp.Compile(`\{.*`)
		helper.IfError(err, "can't (regexp.Compile) to get [jsonExp]")
		ratingJson   := jsonExp.FindAllString(listValue[0], -1)

		if len(ratingJson) == 0 {
			err := json.Unmarshal([]byte(ratingJson[0]), &sRating)
			helper.IfError(err, "can't (json.Unmarshal) to get [sRating = bbbRating{}], maybe Json has been changed") // todo, maybe rebuild error as stateError

			// got json result from service document
			for _, val := range sRating.Search.Results {
				val.BusinessName = strings.ToLower(val.BusinessName)
				grade				 := string([]byte(val.Rating)[0])
	
				// get/check title
				titleText       := strings.ToLower(val.BusinessName)
				titleExp, err   := regexp.Compile("^" + qFeedback.Company + "$")
				helper.IfError(err, "can't (regexp.Compile) to get [titleExp]")
		
				if titleExp.MatchString(titleText) {
					// get/set rate
					if grade == "F" {
						rateText = "1.0"
					} else if grade  == "D" {
						rateText = "2.0"
					} else if grade  == "C" {
						rateText = "3.0"
					} else if grade  == "B" {
						rateText = "4.0"
					} else if grade  == "A" {
						rateText = "5.0"
					} else {
						rateText = "0.0"
						// check if such grade exists
						setParsingErrorByCode(1001, "rate" , errorState)
						continue
						// ------------------------
					}
					foldRate(rateText, sParams, errorState)
				}
			}
		} else {
			setParsingErrorByCode(1005, qFeedback.Company , errorState)
		}
	} else {
		setParsingErrorByCode(1005, qFeedback.Company , errorState)
	}

}
