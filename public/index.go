package public

import (
	"bufio"
	"fmt"
	"gov/backend/interfaces"
	"gov/backend/models"
	"gov/backend/services"
	"errors"
	"strings"
	"os"
)

var feedbackService interfaces.IFeedbackService = &services.FeedbackService{}
var errorService    interfaces.IErrorService    = &services.ErrorService{}

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

		{Title: "tripadvisorua" , Url : "https://www.tripadvisor.com/Search?geo=294473&pid=3826&q=kfc" 										    , ISOCode: "UA", CountryCode: 1  },
		{Title: "otzyvua"       , Url : "https://www.otzyvua.net/search/?q=тонгруп" 									   					       , ISOCode: "UA", CountryCode: 1  },
	}}

	showMsg(console, "")
	
	for scanner.Scan() {
		switch scanner.Text() {
			case "quite" : console.IsQuite = true
			case "q"	 	 : console.IsQuite = true
			default	    : gotKey(console, qfeedback, scanner.Text())
		}
	
		if console.IsQuite {
			break
		} else {
			showMsg(console, scanner.Text())
		}
	}

	errorService.Check(scanner.Err())
}

func gotKey(console *models.ConsoleModel, qfeedback *models.FeedbackQueryModel, text string) {
	feedBacks := make([]models.FeedbackModel, 0)

	if (strings.ToUpper(text) == "UA" || strings.ToUpper(text) == "EU" || strings.ToUpper(text) == "US") && console.Step == 1 {
		qfeedback.ISOCode       = strings.ToUpper(text)
		services, err          := FilterFServiceByISO(qfeedback.Services, qfeedback)
		errorService.Check(err)
		qfeedback.Services      = services
		console.Step            = 2
	} else if strings.ToUpper(text) == "RU" && console.Step == 1 {
		qfeedback.ISOCode       = strings.ToUpper(text)
		services, err          := FilterFServiceByISO(qfeedback.Services, qfeedback)
		errorService.Check(err)
		qfeedback.Services      = services
		console.Step 	         = 3
	} else if console.Step == 1 {
		console.Step 	         = 4
	} else if console.Step == 3 {
		qfeedback.Country       = text
		console.Step 	         = 4
	} else if console.Step == 4 || console.Step == 2 {
		qfeedback.Company       = strings.ToLower(text)
		// Get Average Rate from top Web Portals
		// pass   : feedback.Country, feedback.Company
		// return : feedback.AvarageRate
		
		feedbackService.SetRates(qfeedback, &feedBacks)
		// ----------------

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

		console.IsQuite = true
	}
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
		errorService.Check(errors.New("Unknown Step"))
	}
}

func FilterFServiceByISO(services []*models.FeedbackServiceModel, qfeedback *models.FeedbackQueryModel) ([]*models.FeedbackServiceModel, error) {
	result := make([]*models.FeedbackServiceModel, 0, len(services))
	for _, val := range services {
		 if val.ISOCode == qfeedback.ISOCode {
			result = append(result, val)
		 }
	}

	if len(result) == 0 {
		return nil, errors.New("Reslut is empty")
	} 

	return result, nil
}
