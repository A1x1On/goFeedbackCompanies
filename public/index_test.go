package public

import (
	"testing"
	//"fmt"
	"github.com/stretchr/testify/assert"
	"gov/backend/models"
)

func TestExecInput(t *testing.T) {
	console   := &models.ConsoleModel{IsQuite: false, Step: 1}
	qfeedback := &models.FeedbackQueryModel{Services : []*models.FeedbackServiceModel{
		{Title: "flampRU"         , Url : "https://moscow.flamp.ru/search/{company}" , ISOCode: "RU", CountryCode: 122},
	}}

	qfeedback.ISOCode = "RU"
	enteredText		  := "ru"
	execInput(console, qfeedback, enteredText)

	qfeedback.ISOCode = "US"
	enteredText		   = "Us"
	execInput(console, qfeedback, enteredText)

	qfeedback.ISOCode = "EU"
	enteredText		   = "eU"
	execInput(console, qfeedback, enteredText)

	qfeedback.ISOCode = "UA"
	enteredText		   = "UA"
	execInput(console, qfeedback, enteredText)

	qfeedback.ISOCode = ""
	enteredText		   = ""
	execInput(console, qfeedback, enteredText)

	console.Step	   = 2
	qfeedback.Company = "kfc"
	execInput(console, qfeedback, enteredText)

	console.Step 		= 3
	execInput(console, qfeedback, enteredText)

	console.Step 	   = 4
	execInput(console, qfeedback, enteredText)

	console.Step 	   = 4
	qfeedback.ISOCode = "RU"
	execInput(console, qfeedback, enteredText)
}

func TestShowMsg(t *testing.T) {
	console 		:= &models.ConsoleModel{IsQuite: false, Step: 1}
	enteredText := "US"
	
	showMsg(console, enteredText)
	console.Step = 2
	showMsg(console, enteredText)
	console.Step = 4
	showMsg(console, enteredText)
	//console.Step = 5 // obvious error
	//showMsg(console, enteredText)
}

func TestFilterFServiceByISO(t *testing.T) {
	qfeedback    	:= &models.FeedbackQueryModel{ISOCode : "UA", Services : []*models.FeedbackServiceModel{
		{Title: "flampRU"         , Url : "https://moscow.flamp.ru/search/{company}"											     				         , ISOCode: "RU", CountryCode: 122},
		{Title: "yelpWashington"  , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Washington%2C%20DC"   			   	   , ISOCode: "US", CountryCode: 1  },
		{Title: "yelpPoland"      , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Warszawa%2C%20Mazowieckie%2C%20Poland" 	, ISOCode: "EU", CountryCode: 1  },
		{Title: "otzyvUA"         , Url : "https://www.otzyvua.net/search/?q={company}" 									   					         , ISOCode: "UA", CountryCode: 380},
	}}
	expectedResult := []*models.FeedbackServiceModel{
		{Title: "otzyvUA"         , Url : "https://www.otzyvua.net/search/?q={company}" 									   					         , ISOCode: "UA", CountryCode: 380},
	}
	actualResult	:= filterFServiceByISO(qfeedback.Services, qfeedback)
	assert.Equal(t, expectedResult, actualResult)
}
