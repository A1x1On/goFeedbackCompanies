package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/common/helper"
	"gov/backend/models"
	"math/big"
	"strconv"
	"regexp"
	"math"
)

type serviceParams struct {
	sumRate    float64
	numRate    int
	numReviews int
	doc 		  *goquery.Document
}

func ParseService(doc *goquery.Document, qFeedback *models.FeedbackQueryModel, errorState *models.ErrorStateModel) (float64, int, *models.ErrorStateModel) {
	sParams := &serviceParams{doc: doc}

	// list of services for parsing
	switch qFeedback.Service.Title {
		case "flampRU"	   	  : { ParseFlamp(sParams, errorState, qFeedback)  	} 
		case "yellRU"     	  : { ParseYell(sParams, errorState, qFeedback)    } 
		case "apoiMoscow" 	  : { ParseApoi(sParams, errorState, qFeedback)    } 
		case "pravdaRU"   	  : { ParsePravda(sParams, errorState, qFeedback)	}
		case "spasiboRU"  	  : { ParseSpasibo(sParams, errorState, qFeedback) }
		case "indeedUS"        : { ParseIndeed(sParams, errorState, qFeedback)	}
		case "tripadWashington": { ParseTripad(sParams, errorState, qFeedback)  }
		case "yellowWashington": { ParseYellow(sParams, errorState, qFeedback)  }
		case "bbbUS"			  : { ParseBBB(sParams, errorState, qFeedback)	   }
		case "yelpWashington"  : fallthrough
		case "yelpPoland"      : fallthrough
		case "yelpSpain"       : fallthrough
		case "yelpDenmark"     : fallthrough
		case "yelpBritan"      : fallthrough
		case "yelpNorway"      : { ParseYelp(sParams, errorState, qFeedback)    }
		case "otzyvUA"	        : { ParseOtzyv(sParams, errorState, qFeedback)   }
		default					  : {
			setParsingErrorByCode(1006, qFeedback.ServiceTitle, errorState)
			break
		}
	}

	// check/map types & return data
	if math.IsNaN(sParams.sumRate/float64(sParams.numRate)) {
		return 0, sParams.numReviews, errorState
	} else {
		fixedRate, err	:= strconv.ParseFloat(big.NewFloat(sParams.sumRate/float64(sParams.numRate)).Text('f', 3), 64)
		helper.IfError(err, "can't (trconv.ParseFloat(big.NewFloat(..) to get [fixedRate]")
		return fixedRate, sParams.numReviews, errorState
	}
}


// trim all spaces and new lines
func trimAll(text string) string {
	text = regexp.MustCompile("[\\s\\t\\n]*").ReplaceAllString(text, "")
	return text
}

// get Rate
func foldRate(text string, sParams *serviceParams, errorState *models.ErrorStateModel){
	parsedRate, err := strconv.ParseFloat(trimAll(text), 64)
	helper.IfError(err, "can't (strconv.ParseFloat) to get [parsedRate]")

	// set error if rate has less/more than 0/5 value
	if parsedRate > 5 || parsedRate < 0 {
		setParsingErrorByCode(102, "rate, incorrect val is " + strconv.FormatFloat(parsedRate, 'f', 1, 64), errorState)
	}	
	// ----------------------------------------------

	if parsedRate != 0 {
		sParams.numRate += 1
	}
	sParams.sumRate = sParams.sumRate + parsedRate
}

// get Reviews
func getSumReviews(reviewsText string, numReviews int) int{
	reviewsInt, err := strconv.Atoi(reviewsText)
	helper.IfError(err, "can't (strconv.Atoi) to get [reviewsInt]")
	return numReviews + reviewsInt
}
