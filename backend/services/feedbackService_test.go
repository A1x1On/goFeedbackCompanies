package services

import (
	"testing"
	//"fmt"
	//"github.com/stretchr/testify/assert"
	"gov/backend/models"
)

func TestGetAllByCriteria(t *testing.T) {

	qfeedback := &models.FeedbackQueryModel{
		Company  : "",
		ISOCode	: "",
		Services : []*models.FeedbackServiceModel{
			{Title: "flampRU"         , Url : "https://moscow.flamp.ru/search/{company}"		  			                                              , ISOCode: "RU", CountryCode: 122},
			{Title: "yellRU"          , Url : "https://www.yell.ru/moscow/top/?text={company}"				                                           , ISOCode: "RU", CountryCode: 122},
			{Title: "apoiMoscow"      , Url : "https://www.apoi.ru/kompanii/moskva?searchtext={company}"                                            , ISOCode: "RU", CountryCode: 122},
			{Title: "pravdaRU"        , Url : "https://pravda-sotrudnikov.ru/search?q={company}"			                                           , ISOCode: "RU", CountryCode: 122},
			{Title: "spasiboRU"       , Url : "https://spasibovsem.ru/search/?q={company}" 					                                           , ISOCode: "RU", CountryCode: 122},

			{Title: "indeedUS"        , Url : "https://www.indeed.com/cmp?q={company}&l=&from=discovery-cmp-search"     					    	 		 , ISOCode: "US", CountryCode: 1  },
			{Title: "yelpWashington"  , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Washington%2C%20DC"   			   	 		 , ISOCode: "US", CountryCode: 1  },
			{Title: "tripadWashington", Url : "https://www.yellowpages.com/search?search_terms={company}&geo_location_terms=Washington%2C+DC" 		 , ISOCode: "US", CountryCode: 1  },
			{Title: "bbbUS"           , Url : "https://www.bbb.org/search?filter_ratings=F&find_country=USA&find_text={company}&page=1&sort=Rating" , ISOCode: "US", CountryCode: 1  },
			{Title: "yellowWashington", Url : "https://www.yellowpages.com/search?search_terms={company}&geo_location_terms=Washington%2C+DC"		 , ISOCode: "US", CountryCode: 1  },

			{Title: "yelpBritan"      , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=United%20Kingdom%20London" 		  				 , ISOCode: "EU", CountryCode: 1  },
			{Title: "yelpNorway"      , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Norway%20Oslo" 		  							    , ISOCode: "EU", CountryCode: 1  },
			{Title: "yelpPoland"      , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Warszawa%2C%20Mazowieckie%2C%20Poland" 	    , ISOCode: "EU", CountryCode: 1  },
			{Title: "yelpSpain"       , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Madrid%2C%20Spain" 		  					    , ISOCode: "EU", CountryCode: 1  },
			{Title: "yelpDenmark"     , Url : "https://www.yelp.com/search?find_desc={company}&find_loc=Denmark%20Copenhagen"         				    , ISOCode: "EU", CountryCode: 1  },

			{Title: "otzyvUA"         , Url : "https://www.otzyvua.net/search/?q={company}" 									   					       		 , ISOCode: "UA", CountryCode: 380},
		},
	}

	feedbackService := FeedbackService{}
	actualResult    := feedbackService.GetAllByCriteria(qfeedback)

	if len(actualResult) == 0 {
		t.Error("Services is Empty")
	}
}