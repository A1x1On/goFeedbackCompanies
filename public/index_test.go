package public

import (
	"testing"
	"fmt"

	//"github.com/stretchr/testify/assert"
	"gov/backend/models"
)

func TestGotKey(t *testing.T) {
	fmt.Println("***************************TestGotKey*************************************")

	console   := &models.ConsoleModel{IsQuite: false, Step: 1}
	qfeedback := &models.FeedbackQueryModel{Services : []*models.FeedbackServiceModel{
		{Title: "flamp", Url : "https://{country}.flamp.ru/search/{company}", ISOCode: "RU", CountryCode: 122},
	}}

	qfeedback.ISOCode = "RU"
	gotKey(console, qfeedback, "RU")

	qfeedback.ISOCode = "US"
	gotKey(console, qfeedback, "US")

	qfeedback.ISOCode = "EU"
	gotKey(console, qfeedback, "EU")

	qfeedback.ISOCode = "UA"
	gotKey(console, qfeedback, "UA")

	qfeedback.ISOCode = ""
	gotKey(console, qfeedback, "")
	console.Step	   = 2
	qfeedback.Company = "kfc"
	gotKey(console, qfeedback, "")

	console.Step 		= 3
	qfeedback.Country = "moscow"
	gotKey(console, qfeedback, "")

	console.Step 	   = 4
	gotKey(console, qfeedback, "")

	console.Step 	   = 4
	qfeedback.ISOCode = "RU"
	qfeedback.Country = ""
	gotKey(console, qfeedback, "RU")


	//assert.Equal(t, expectedResult, actualResult)
}

func TestShowMsg(t *testing.T) {
	fmt.Println("***************************TestShowMsg*************************************")

	console   := &models.ConsoleModel{IsQuite: false, Step: 1}
	
	showMsg(console, "US")
	console.Step = 2
	showMsg(console, "US")
	console.Step = 3
	showMsg(console, "US")
	console.Step = 4
	showMsg(console, "US")
	//console.Step = 5 // ERROR
	//showMsg(console, "US")

}