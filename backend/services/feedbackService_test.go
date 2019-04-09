package services

import (
	"testing"
	//"fmt"
	//"github.com/stretchr/testify/assert"
	"gov/backend/models"
)

func TestGetAllByCriteria(t *testing.T) {

	qfeedback := &models.FeedbackQueryModel{
		Country	: "",
		Company  : "",
		ISOCode	: "",
		Services : []*models.FeedbackServiceModel{
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
		},
	}

	feedbackService := FeedbackService{}
	actualResult    := feedbackService.GetAllByCriteria(qfeedback)

	if len(actualResult) == 0 {
		t.Error("Services is Empty")
	}
}