package services

import (
	"gov/backend/services/auxiliary/feedBackAux"
	"gov/backend/common/connections"
	"gov/backend/common/helper"
	"github.com/ztrue/tracerr"
	"gov/backend/common/config"
	"gov/backend/repositories"
	"gov/backend/interfaces"
	"gov/backend/models"
	"math/big"
	"strconv"
	"regexp"
	"fmt"
)

var feedbackRepository interfaces.IFeedbackRepository = &repositories.FeedbackRepository{}

type FeedbackService struct{}

func (s *FeedbackService) GetAll() string{
	//params := map[string]string{"login" : "2tzEhq13","pass" : "bHXAG7sJ",}
	response := connections.ProxyRequest("GET",  config.Set.API.BaseURL.Golang + "x/net/prox", nil)
	return string(response)
}

var services = []models.FeedbackServiceModel{
	{Id: 1,  Title: "flampRU"         , Url : "https://moscow.flamp.ru/search/{company}"		  			                                              , ISOCode: "RU", CountryCode: 122},
	{Id: 2,  Title: "yellRU"          , Url : "https://www.yell.ru/moscow/top/?text={company}"				                                           , ISOCode: "RU", CountryCode: 122},
	{Id: 3,  Title: "apoiMoscow"      , Url : "https://www.apoi.ru/kompanii/moskva?searchtext={company}"                                             , ISOCode: "RU", CountryCode: 122},
	{Id: 4,  Title: "pravdaRU"        , Url : "https://pravda-sotrudnikov.ru/search?q={company}"			                                              , ISOCode: "RU", CountryCode: 122},
	{Id: 5,  Title: "spasiboRU"       , Url : "https://spasibovsem.ru/search/?q={company}" 					                                           , ISOCode: "RU", CountryCode: 122},

	{Id: 6,  Title: "indeedUS"        , Url : "https://www.indeed.com/cmp?q={company}&l=&from=discovery-cmp-search"     					    	 		    , ISOCode: "US", CountryCode: 1  },
	{Id: 7,  Title: "yelpWashington"  , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Washington%2C%20DC"   			   	 		    , ISOCode: "US", CountryCode: 1  },
	{Id: 8,  Title: "tripadWashington", Url : "https://www.yellowpages.com/search?search_terms={company}&geo_location_terms=Washington%2C+DC" 		 , ISOCode: "US", CountryCode: 1  },
	{Id: 9,  Title: "bbbUS"           , Url : "https://www.bbb.org/search?filter_ratings=F&find_country=USA&find_text={company}&page=1&sort=Rating"  , ISOCode: "US", CountryCode: 1  },
	{Id: 10, Title: "yellowWashington", Url : "https://www.yellowpages.com/search?search_terms={company}&geo_location_terms=Washington%2C+DC"		    , ISOCode: "US", CountryCode: 1  },

	{Id: 11, Title: "yelpBritan"      , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=United%20Kingdom%20London" 		  				 , ISOCode: "EU", CountryCode: 100},
	{Id: 12, Title: "yelpNorway"      , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Norway%20Oslo" 		  							    , ISOCode: "EU", CountryCode: 100},
	{Id: 13, Title: "yelpPoland"      , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Warszawa%2C%20Mazowieckie%2C%20Poland" 	    , ISOCode: "EU", CountryCode: 100},
	{Id: 14, Title: "yelpSpain"       , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Madrid%2C%20Spain" 		  					       , ISOCode: "EU", CountryCode: 100},
	{Id: 15, Title: "yelpDenmark"     , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Denmark%20Copenhagen"         				    , ISOCode: "EU", CountryCode: 100},

	{Id: 16, Title: "otzyvUA"         , Url : "https://www.otzyvua.net/search/?q={company}" 									   					       		 , ISOCode: "UA", CountryCode: 380},
}

func (s *FeedbackService) GetReviewService(feedbackParams *models.FeedbcakParamsModel) (models.FeedbackModel, int, string) {
	qFeedback  := &models.FeedbackQueryModel{Company : feedbackParams.Company,}
	for _, val := range services {
		if val.Id == feedbackParams.ServiceId {
			qFeedback.Service = val
			continue
		}
	}

	if qFeedback.Service.Id == 0 {
		helper.IfError(tracerr.New("Service is not found on (GetReviewService func)"), "Service with [" + strconv.Itoa(feedbackParams.ServiceId) + "] id is not found ")
	}

	fmt.Println("'" + qFeedback.Service.Title + "' is connecting ...", )

	qFeedback.Service.Url  = getRealUrl(qFeedback) // replace templates keys into the data
	doc, code, err        := feedbackRepository.GetFeedbackPage(qFeedback.Service.Url)       // get page of the service
	errorState            := &models.ErrorStateModel{Message: "null", Code: 0}		 // init new errorState

	if code != 200 || err != nil {
		if err != nil {
			errorState = &models.ErrorStateModel{Message : err.Error(), Code : code,}
		} else {
			code = feedBackAux.VerifyNotFoundPage(doc, qFeedback, errorState) // some services can return 404 but it can mean just not found results, that checks
			feedBackAux.SetHttpErrorByCode(code, errorState)
		}

		return models.FeedbackModel{ServiceTitle: qFeedback.Service.Title, Rate: 0.0, NumReviews: 0, ErrorState: errorState}, errorState.Code, errorState.Message
	}

	// parse got html for passed current service & get MAIN feedback data
	rate, numReviews, errState := feedBackAux.ParseService(doc, qFeedback, errorState)
	// ------------------------------------------------------------------

	return models.FeedbackModel{ServiceTitle: qFeedback.Service.Title, Rate: rate, NumReviews: numReviews, ErrorState: errState}, errState.Code, errState.Message
}

func (s *FeedbackService) GetServices(CountryCode int) ([]int, int, string) {
	errorCode  := 0
	errorMsg   := ""
	result 	  := make([]int, 0)
	for _, val := range services {
		 if val.CountryCode == CountryCode {
			result = append(result, val.Id)
		 }
	}

	if len(result) == 0 {
		if CountryCode == 0 {
			errorCode = 404
			errorMsg  = "CountryCode is 0, result is empty"
		} else {
			errorCode = 404
			errorMsg  = "Services have not been found by current CountryCode"
		}
	}

	return result, errorCode, errorMsg
}


func (s *FeedbackService) GetAllByCriteria(qFeedback *models.FeedbackQueryModel) []models.FeedbackModel{
	feedBacks 				:= make([]models.FeedbackModel, 0)
	averageBy            := 0.0
	totalFeedback        := len(qFeedback.Services)

	// exec and append ready feedback for each service
	for _, service := range qFeedback.Services {
		fmt.Println("'" + service.Title + "' is connecting ...", )

		qFeedback.ServiceTitle = service.Title
		service.Url            = getReplacedUrl(qFeedback, service.Url, service.Title) // replace templates keys into the data
		doc, code, err        := feedbackRepository.GetFeedbackPage(service.Url)       // get page of the service
		errorState            := &models.ErrorStateModel{Message: "null", Code: 0}		 // init new errorState

		if code != 200 || err != nil {
			if err != nil {
				errorState = &models.ErrorStateModel{Message : err.Error(), Code : code,}
			} else {
				code = feedBackAux.VerifyNotFoundPage(doc, qFeedback, errorState) // some services can return 404 but it can mean just not found results, that checks
				feedBackAux.SetHttpErrorByCode(code, errorState)
			}

			feedBacks = append(feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: 0.0, NumReviews: 0, ErrorState: errorState})
			continue
		}

		// parse got html for passed current service & get MAIN feedback data
		rate, numReviews, errState := feedBackAux.ParseService(doc, qFeedback, errorState)
		// ------------------------------------------------------------------
		averageBy 				       = averageBy + rate
		qFeedback.NumReviews        = qFeedback.NumReviews + numReviews
		feedBacks 			          = append(feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: rate, NumReviews: numReviews, ErrorState: errState})

		if rate == 0 {
			totalFeedback = totalFeedback - 1
		}
	}

	if(averageBy == 0){
		qFeedback.AvarageRate = 0
	} else {
		qFeedback.AvarageRate = averageBy/float64(totalFeedback)
		fixedRate, err			 := strconv.ParseFloat(big.NewFloat(averageBy/float64(totalFeedback)).Text('f', 3), 64)
		helper.IfError(err, "can't (strconv.ParseFloat) to get [fixedRate]")
		qFeedback.AvarageRate  = fixedRate
	}

	return feedBacks
}

func getReplacedUrl(qfeedback *models.FeedbackQueryModel, url string, title string) string{
	return regexp.MustCompile("{company}").ReplaceAllString(url, regexp.MustCompile("\\s").ReplaceAllString(qfeedback.Company, "+"))
}

func getRealUrl(qfeedback *models.FeedbackQueryModel) string{
	return regexp.MustCompile("{company}").ReplaceAllString(qfeedback.Service.Url, regexp.MustCompile("\\s").ReplaceAllString(qfeedback.Company, "+"))
}


