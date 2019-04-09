package public

import (
	"bufio"
	"fmt"
	"gov/backend/interfaces"
	"gov/backend/models"
	"gov/backend/services"
	"gov/backend/common/helpers"
	"errors"
	"strings"
	"os"
)

var feedbackService interfaces.IFeedbackService = &services.FeedbackService{}

func Index() {
	scanner   := bufio.NewScanner(os.Stdin)
	console   := &models.ConsoleModel{IsQuite: false, Step: 1}
	qfeedback := &models.FeedbackQueryModel{Services : []*models.FeedbackServiceModel{
		{Title: "flamp"         , Url : "https://{country}.flamp.ru/search/{company}"		  			                                     , ISOCode: "RU", CountryCode: 122},
		{Title: "yell"          , Url : "https://www.yell.ru/{country}/top/?text={company}"				                                  , ISOCode: "RU", CountryCode: 122},
		{Title: "apoi"          , Url : "https://www.apoi.ru/kompanii/{country}?searchtext={company}"                                  , ISOCode: "RU", CountryCode: 122},
		{Title: "pravda"        , Url : "https://pravda-sotrudnikov.ru/search?q={company}"			                                     , ISOCode: "RU", CountryCode: 122},
		{Title: "spasibo"       , Url : "https://spasibovsem.ru/search/?q={company}" 					                                     , ISOCode: "RU", CountryCode: 122},

		{Title: "indeedcom"     , Url : "https://www.indeed.com/cmp?q={company}&l=&from=discovery-cmp-search"     					    	 , ISOCode: "US", CountryCode: 1  },
		{Title: "yelp"          , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Washington%2C%20DC"   			   	 , ISOCode: "US", CountryCode: 1  },
		{Title: "tripadvisorus" , Url : "https://www.yellowpages.com/search?search_terms={company}&geo_location_terms=Washington%2C+DC", ISOCode: "US", CountryCode: 1  },
		{Title: "bbb"           , Url : "https://www.bbb.org/search?find_country=USA&find_text=kfc&page=1"  							       , ISOCode: "US", CountryCode: 1  },
		{Title: "yellowpages"   , Url : "https://www.yellowpages.com/search?search_terms={company}&geo_location_terms=Washington%2C+DC", ISOCode: "US", CountryCode: 1  },

		{Title: "tripadvisoreu" , Url : "https://www.tripadvisor.com/Search?geo=4&pid=3826&q=kfc" 		  							          , ISOCode: "EU", CountryCode: 1  },
		{Title: "pagesjaunes"   , Url : "https://www.pagesjaunes.fr/recherche/berling-57/kfc" 		  							                , ISOCode: "EU", CountryCode: 1  },
		{Title: "yelpPoland"    , Url : "https://www.yelp.com/search?find_desc=kfc&find_loc=Warszawa%2C%20Mazowieckie%2C%20Poland" 	 , ISOCode: "EU", CountryCode: 1  },
		{Title: "yelpSpain"     , Url : "https://www.yelp.com/search?find_desc=kfc&find_loc=Madrid%2C%20Spain" 		  					    , ISOCode: "EU", CountryCode: 1  },
		{Title: "yelpDenmark"   , Url : "https://www.yelp.com/search?find_desc=kfc&find_loc=Copenhagen%2C%20Hovedstaden%2C%20Denmark"  , ISOCode: "EU", CountryCode: 1  },

		{Title: "tripadvisorua" , Url : "https://www.tripadvisor.com/Search?geo=294473&pid=3826&q=kfc" 										    , ISOCode: "UA", CountryCode: 380},
		{Title: "otzyvua"       , Url : "https://www.otzyvua.net/search/?q=тонгруп" 									   					       , ISOCode: "UA", CountryCode: 380},
	}}

	showMsg(console, "") // display first text instruction into the console
	
	// begin keyboard listening ... 
	for scanner.Scan() {
		switch scanner.Text() {
			case "quite" : console.IsQuite = true
			case "q"	 	 : console.IsQuite = true
			default	    : execInput(console, qfeedback, scanner.Text()) // pick the appropriate console step in the condition blocks
		}
	
		if console.IsQuite { // if IsQuite == true do console exit
			break
		} else {
			showMsg(console, scanner.Text()) // display text instruction for next step
		}
	}

	helpers.CheckError(scanner.Err())
}

func showMsg(console *models.ConsoleModel, text string) {
	switch console.Step {
	case 1:
		fmt.Println("-----------------------\n|FEEDBACK APP IS READY|\n-----------------------\nEnter ISO Codes, please:\n(option: you can enter 'RU', 'UA', 'EU' or 'US' for strict ISO Code search)\nEnter Empty line to set All")
	case 2:
		fmt.Println("SET '" + strings.ToUpper(text) + "' ISO Code")	
		fmt.Println("Enter Company, please: ")
	case 3:
		fmt.Println("SET 'RU' ISO Code")	
		fmt.Println("(option: you can enter 'Country' for strict search)\nEnter Empty line to set All")
	case 4:
		fmt.Println("Enter Company, please: ")
	default:
		helpers.CheckError(errors.New("Unknown Step"))
	}
}

func execInput(console *models.ConsoleModel, qfeedback *models.FeedbackQueryModel, textKey string) {
	feedBacks := make([]models.FeedbackModel, 0)
	textKey    = strings.ToUpper(textKey)

	if (textKey== "UA" || textKey == "EU" || textKey == "US") && console.Step == 1 { // if is other zone
		qfeedback.ISOCode       = textKey
		services	              := filterFServiceByISO(qfeedback.Services, qfeedback) // filter/get services by entered ISOCode
		qfeedback.Services      = services
		console.Step            = 2
	} else if textKey == "RU" && console.Step == 1 { // if is RU zone
		qfeedback.ISOCode       = textKey
		services               := filterFServiceByISO(qfeedback.Services, qfeedback)
		qfeedback.Services      = services
		console.Step 	         = 3
	} else if console.Step == 1 { // if are All available zones
		console.Step 	         = 4
	} else if console.Step == 3 { // if is country
		qfeedback.Country       = strings.ToLower(textKey)
		console.Step 	         = 4
	} else if (console.Step == 4 || console.Step == 2) && textKey != "" { // if is company
		qfeedback.Company       = strings.ToLower(textKey)
		
		// Set Average Rate & Count of Reviews for All services
		// Get Struct Array about Found feedback services
		feedBacks = feedbackService.GetAllByCriteria(qfeedback)
		// ----------------------------------------------------

		fmt.Println("============== Feedbacks have been prepared =================")
		for _, feedback := range feedBacks {
			fmt.Println("-------------------------------------")
			fmt.Println("Service Title: "  , feedback.ServiceTitle)
			fmt.Println("Average Rate: "   , feedback.Rate)
			fmt.Println("Review Count: "   , feedback.NumReviews)
			fmt.Println("State of Result: ", feedback.StateResult)
		}

		fmt.Println("========== In Total of the Services ================")
		fmt.Println("Average Rate: ", qfeedback.AvarageRate)
		fmt.Println("Review Count: ", qfeedback.NumReviews)

		console.IsQuite = true // console exit
	}
}

func filterFServiceByISO(services []*models.FeedbackServiceModel, qfeedback *models.FeedbackQueryModel) []*models.FeedbackServiceModel {
	result := make([]*models.FeedbackServiceModel, 0, len(services))

	fmt.Println("----------- services[0].ISOCode----------", services[0].ISOCode)
	for _, val := range services {
		 if val.ISOCode == qfeedback.ISOCode {
			result = append(result, val)
		 }
	}

	if len(result) == 0 {
		if qfeedback.ISOCode == "" {
			helpers.CheckError(errors.New("qfeedback.ISOCode is empty => Reslut is empty"))
		} else {
			helpers.CheckError(errors.New("Reslut is empty"))
		}
	} 

	return result
}
