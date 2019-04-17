package services

import (
	"gov/backend/services/auxiliary/feedBackAux"
	"gov/backend/common/connections"
	"gov/backend/common/helper"
	"gov/backend/common/config"
	"gov/backend/repositories"
	"gov/backend/interfaces"
	"github.com/pkg/errors"
	"gov/backend/models"
	"math/big"
	"strconv"
	"regexp"
	"fmt"	
)

var feedbackRepository interfaces.IFeedbackRepository = &repositories.FeedbackRepository{}

type FeedbackService struct{}

func (s *FeedbackService) GetAll() string{
	baseUrl := config.Set.API.BaseURL.Golang
	//params := map[string]string{"login" : "2tzEhq13","pass" : "bHXAG7sJ",}
	response := connections.ProxyRequest("GET",  baseUrl + "x/net/prox", nil)
	return string(response)
}

func (s *FeedbackService) GetAllByCriteria(qFeedback *models.FeedbackQueryModel) []models.FeedbackModel{
	feedBacks 				:= make([]models.FeedbackModel, 0)
	errorState 				:= &models.ErrorStateModel{}
	qFeedback.AvarageRate = 0.0
	averageBy            := 0.0
	totalFeedback        := len(qFeedback.Services)
	urlDataSet	         := func(qfeedback *models.FeedbackQueryModel, sFeedback *models.FeedbackServiceModel){		
		if sFeedback.Title == "apoiMoscow" {
			if qfeedback.Country == "" || qfeedback.Country == "moscow" {
				qfeedback.Country = "moskva"
			}
		} else if qfeedback.Country == "" {
			qfeedback.Country = "moscow"
		}

		sFeedback.Url = regexp.MustCompile("{country}").ReplaceAllString(sFeedback.Url, qfeedback.Country)
		sFeedback.Url = regexp.MustCompile("{company}").ReplaceAllString(sFeedback.Url, regexp.MustCompile("\\s").ReplaceAllString(qfeedback.Company, "+"))
	}

	// exec and append ready feedback for each service
	for _, service := range qFeedback.Services {
		if qFeedback.Country != "" && (service.Title == "pravda" || service.Title == "spasibo") {
			totalFeedback = totalFeedback - 1
			continue
		}
		
		urlDataSet(qFeedback, service) // replace templates keys into the data
		fmt.Println("'" + service.Title + "' is connecting ...", )
		
		doc, err := feedbackRepository.GetFeedbackPage(service.Url) // get page of the service
		if err != nil {
			errorState = &models.ErrorStateModel{
				Message : err.Error(),
				Code    : 666,
			}

			feedBacks = append(feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: 0.0, NumReviews: 0, ErrorState: errorState})
			continue
		}

		// parse got html for passed current service & get MAIN feedback data
		rate, numReviews, errState := feedBackAux.ParseService(doc, qFeedback, service.Title)
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
		helper.CheckError(errors.Wrap(err, "can't (strconv.ParseFloat) to get [fixedRate]"))
		qFeedback.AvarageRate  = fixedRate
	}

	return feedBacks
}


