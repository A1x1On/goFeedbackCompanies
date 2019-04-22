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

func ParseService(doc *goquery.Document, qFeedback *models.FeedbackQueryModel, title string) (float64, int, *models.ErrorStateModel) {
	errorState := &models.ErrorStateModel{Message: "null", Code: 0}
	sParams	  := &serviceParams{doc: doc}

	// list of services for parsing
	switch title {
		case "flampRU"	   	  : { ParseFlamp(sParams, errorState)  } // todo add more parsing errors
		case "yellRU"     	  : { ParseYell(sParams, errorState, qFeedback)    } // todo add more parsing errors
		case "apoiMoscow" 	  : { ParseApoi(sParams, errorState, qFeedback)    } 
		case "pravdaRU"   	  : { ParsePravda(sParams, errorState, qFeedback)	} // todo add more parsing errors
		case "spasiboRU"  	  : { ParseSpasibo(sParams, errorState, qFeedback) } // todo add more parsing errors
		case "indeedUS"        : { ParseIndeed(sParams, errorState, qFeedback)	} // todo add more parsing errors
		case "tripadWashington": { ParseTripad(sParams, errorState, qFeedback)  } // todo add more parsing errors
		case "yellowWashington": { ParseYellow(sParams, errorState, qFeedback)  } // todo add more parsing errors
		case "bbbUS"			  : { ParseBBB(sParams, errorState, qFeedback)	   } // todo add more parsing errors
		case "yelpWashington"  : { ParseYelp(sParams, errorState, qFeedback)    } // todo add more parsing errors
		case "yelpPoland"      : { ParseYelp(sParams, errorState, qFeedback)    } // todo add more parsing errors
		case "yelpSpain"       : { ParseYelp(sParams, errorState, qFeedback)    } // todo add more parsing errors
		case "yelpDenmark"     : { ParseYelp(sParams, errorState, qFeedback)    } // todo add more parsing errors
		case "yelpBritan"      : { ParseYelp(sParams, errorState, qFeedback)    } // todo add more parsing errors
		case "yelpNorway"      : { ParseYelp(sParams, errorState, qFeedback)    } // todo add more parsing errors
		case "otzyvUA"	        : { ParseOtzyv(sParams, errorState, qFeedback)    } // todo add more parsing errors
		default					  : {
			setHttpErrorByCode(404, errorState)
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

// Aux funcs ---------------------------------------
// trim all spaces and new lines
func trimAll(text string) string {
	text = regexp.MustCompile("[\\s\\t\\n]*").ReplaceAllString(text, "")
	return text
}

// get Rate
func foldRate(text string, sParams *serviceParams){
	parsedRate, err := strconv.ParseFloat(trimAll(text), 64)
	helper.IfError(err, "can't (strconv.ParseFloat) to get [parsedRate]")
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
// -------------------------------------------------

