package services

import (
	"fmt"
	//"github.com/PuerkitoBio/goquery"
	"gov/backend/interfaces"
	"gov/backend/models"
	"gov/backend/repositories"
	"gov/backend/services/auxiliary/feedBackAux"
	//"gov/backend/common/helpers"
	//"strings"
	//"math"
	"regexp"
	//"strconv"
)

var feedbackRepository interfaces.IFeedbackRepository = &repositories.FeedbackRepository{}

type FeedbackService struct{}

func (s *FeedbackService) GetAllByCriteria(qFeedback *models.FeedbackQueryModel) []models.FeedbackModel{
	feedBacks 				:= make([]models.FeedbackModel, 0)
	qFeedback.AvarageRate = 0.0
	avarageBy            := 0.0
	totalFeedback        := len(qFeedback.Services)
	urlDataSet	         := func(qfeedback *models.FeedbackQueryModel, sFeedback *models.FeedbackServiceModel){		
		if sFeedback.Title == "apoi" {
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

		// replace templates keys into the data
		urlDataSet(qFeedback, service)
		fmt.Println("'" + service.Title + "' is connecting ...", )
		// get page of the service
		doc, err := feedbackRepository.GetFeedbackPage(service.Url)
		if err != nil {
			// add feedback error data
			feedBacks = append(feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: 0.0, NumReviews: 0, StateResult: err.Error()})
			continue
		}

		// get main feedback data
		rate, numReviews, state := feedBackAux.ParseService(doc, qFeedback, service.Title)
		avarageBy 				    = avarageBy + rate
		qFeedback.NumReviews     = qFeedback.NumReviews + numReviews
		// add main feedback data
		feedBacks 			       = append(feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: rate, NumReviews: numReviews, StateResult: state})

		if rate == 0 {
			totalFeedback = totalFeedback - 1
		}
	}

	if(avarageBy == 0){
		qFeedback.AvarageRate = 0
	} else {
		qFeedback.AvarageRate = avarageBy/float64(totalFeedback) 
	}

	return feedBacks
}


