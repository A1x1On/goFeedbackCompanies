package services

import (
	"fmt"
	"gov/backend/interfaces"
	"gov/backend/models"
	"gov/backend/repositories"
	"gov/backend/common/helpers"
	"gov/backend/services/auxiliary/feedBackAux"
	"regexp"
	"math/big"
	"strconv"
)

var feedbackRepository interfaces.IFeedbackRepository = &repositories.FeedbackRepository{}

type FeedbackService struct{}

func (s *FeedbackService) GetAllByCriteria(qFeedback *models.FeedbackQueryModel) []models.FeedbackModel{
	feedBacks 				:= make([]models.FeedbackModel, 0)
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
			feedBacks = append(feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: 0.0, NumReviews: 0, StateResult: err.Error()})
			continue
		}

		// parse got html for passed current service & get MAIN feedback data
		rate, numReviews, state := feedBackAux.ParseService(doc, qFeedback, service.Title)
		// ------------------------------------------------------------------
		averageBy 				    = averageBy + rate
		qFeedback.NumReviews     = qFeedback.NumReviews + numReviews
		feedBacks 			       = append(feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: rate, NumReviews: numReviews, StateResult: state})

		if rate == 0 {
			totalFeedback = totalFeedback - 1
		}
	}

	if(averageBy == 0){
		qFeedback.AvarageRate = 0
	} else {
		qFeedback.AvarageRate = averageBy/float64(totalFeedback) 
		fixedRate, err			 := strconv.ParseFloat(big.NewFloat(averageBy/float64(totalFeedback)).Text('f', 3), 64)
		helpers.CheckError(err)
		qFeedback.AvarageRate  = fixedRate
	}

	return feedBacks
}


