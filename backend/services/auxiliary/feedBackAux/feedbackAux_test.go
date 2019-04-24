package feedBackAux

import (
	"github.com/stretchr/testify/assert"
	"gov/backend/repositories"
	"gov/backend/interfaces"
	"github.com/pkg/errors"
	"gov/backend/models"
	"testing"
	"strconv"
)

var feedbackRepository interfaces.IFeedbackRepository = &repositories.FeedbackRepository{}

func TestParseService(t *testing.T) {
	qfeedback    	:= &models.FeedbackQueryModel{ServiceTitle : "yelpPoland", Services : []*models.FeedbackServiceModel{
		{Title: "flamp"         , Url : "https://{country}.flamp.ru/search/{company}"											     				, ISOCode: "RU", CountryCode: 122},
		{Title: "yelp"          , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Washington%2C%20DC"    				, ISOCode: "US", CountryCode: 1  },
		{Title: "yelpPoland"    , Url : "https://www.yelp.com/search?find_desc=kfc&find_loc=Warszawa%2C%20Mazowieckie%2C%20Poland" , ISOCode: "EU", CountryCode: 1  },
		{Title: "tripadvisorua" , Url : "https://www.tripadvisor.com/Search?geo=294473&pid=3826&q=kfc" 									   , ISOCode: "UA", CountryCode: 380},
	}}

	doc, code, err := feedbackRepository.GetFeedbackPage("https://www.tripadvisor.com/Search?geo=4&pid=3826&q=kfc") // get page of the service

	if code != 200 || err != nil {
		if err != nil {
			assert.Error(t, err)
		} else {
			strCode := strconv.Itoa(code)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Error(t, errors.New(strCode))
			}
		}
	} else {
		errorState                 := &models.ErrorStateModel{Message: "null", Code: 0}
		rate, numReviews, errState := ParseService(doc, qfeedback, errorState)

		var expectedRate 		  float64
		var expectedNumReviews int
		var expectedstate      *models.ErrorStateModel

		assert.IsType(t, rate		 , expectedRate)
		assert.IsType(t, numReviews , expectedNumReviews)
		assert.IsType(t, errState   , expectedstate)
	}
}