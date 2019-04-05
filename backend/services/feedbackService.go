package services

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gov/backend/interfaces"
	"gov/backend/models"
	"gov/backend/repositories"
	"strings"
	"math"
	"regexp"
	"strconv"
)

type FeedbackService struct{}

var feedbackRepository interfaces.IFeedbackRepository = &repositories.FeedbackRepository{}
var errorService       interfaces.IErrorService       = &ErrorService{}

func (s *FeedbackService) SetRates(qFeedback *models.FeedbackQueryModel, feedBacks *[]models.FeedbackModel){
	qFeedback.AvarageRate = 0.0
	avarageBy            := 0.0
	totalFeedback        := len(qFeedback.Services)
	urlDataSet	         := func(qfeedback *models.FeedbackQueryModel, sFeedback *models.FeedbackServiceModel){		
		if qfeedback.Country == "" {
			countryStr := ""
			if sFeedback.Title == "apoi" {
				countryStr = "moskva"
			} else {
				countryStr = "moscow"
			}
			sFeedback.Url = regexp.MustCompile("{country}").ReplaceAllString(sFeedback.Url, countryStr)
		}else{
			sFeedback.Url = regexp.MustCompile("{country}").ReplaceAllString(sFeedback.Url, qfeedback.Country)
		}
		sFeedback.Url = regexp.MustCompile("{company}").ReplaceAllString(sFeedback.Url, regexp.MustCompile("\\s").ReplaceAllString(qfeedback.Company, "+"))
	}

	// exec and append ready feedback for each service
	for _, service := range qFeedback.Services {
		if qFeedback.Country != "" && (service.Title == "pravda" || service.Title == "spasibo") {
			totalFeedback = totalFeedback - 1
			continue
		}

		// replace templates key into the data
		urlDataSet(qFeedback, service)
		// get page of service
		doc, err := feedbackRepository.GetFeedbackPage(service.Url)
		if err != nil {
			// add feedback error data
			*feedBacks = append(*feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: 0.0, NumReviews: 0, StateResult: err.Error()})
			continue
		}

		// get main feedback data
		rate, numReviews, state := parseService(doc, qFeedback, service.Title)
		avarageBy 				    = avarageBy + rate
		qFeedback.NumReviews     = qFeedback.NumReviews + numReviews
		// add main feedback data
		*feedBacks 			       = append(*feedBacks, models.FeedbackModel{ServiceTitle: service.Title, Rate: rate, NumReviews: numReviews, StateResult: state})

		if rate == 0 {
			totalFeedback = totalFeedback - 1
		}
	}

	if(avarageBy == 0){
		qFeedback.AvarageRate = 0
	} else {
		qFeedback.AvarageRate = avarageBy/float64(totalFeedback) 
	}
}

func parseService(doc *goquery.Document, qFeedback *models.FeedbackQueryModel, title string) (float64, int, string) {
	state		     := "Unknown Status"
	sumRate		  := 0.0
	numRate       := 0
	numReviews    := 0
	setSumRate    := func(text string){
		parsedRate, err := strconv.ParseFloat(trim(text), 64)
		errorService.Check(err)
		if parsedRate != 0 {
			sumRate    = sumRate + parsedRate
			numRate += 1
		}
	}
	setSumReviews := func(reviewsText string) int{
		reviewsInt, _ := strconv.Atoi(reviewsText)
		return numReviews + reviewsInt
	}

	fmt.Println("'" + title + "' is parsing ...", )
	// take one
	switch title {
		case "flamp": {
			doc.Find("cat-brand-filial-rating").Each(func(i int, sel *goquery.Selection) {
				// get/set rate
				rateText, _          := sel.Attr("rating")
				setSumRate(rateText)
				// get/set reviews
				reviewsText, _       := sel.Attr("reviews-count")
				numReviews	   		 = setSumReviews(reviewsText)
			})

			state = "It's Parsed Successfully"
		} 
		case "yell": {
			doc.Find("div.companies__item-content").Each(func(i int, sel *goquery.Selection) {	
				// get/set title
				titleHtml     := sel.Find(".companies__item-title-text")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
				errorService.Check(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml 	  := sel.Find("span.rating__value")
					rateText 	  := trim(strings.ToLower(rateHtml.Text()))
					setSumRate(rateText)
					// get/set reviews
					reviewsHtml   := sel.Find("span.rating__reviews > span")
					numReviews	   = setSumReviews(reviewsHtml.Text())
				}
			})

			state = "It's Parsed Successfully"
		} 
		case "apoi": {
			doc.Find("div.w_star").Each(func(i int, sel *goquery.Selection) {
				// get/set title
				titleHtml     := sel.Find(".m_title > .flw a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
				errorService.Check(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml      := sel.Find(".img_p > .p_f_s")
					rateExp, err  := regexp.Compile("[\\d\\.]*")
					rateText      := rateExp.FindAllString(trim(rateHtml.Text()), -1)
					errorService.Check(err)
					setSumRate(rateText[0])
					// get/set reviews
					reviewsHtml   := sel.Find(".img_p > .numReviews")
					reviewsText	  := regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
					numReviews	   = setSumReviews(reviewsText)
				}
			})

			state = "It's Parsed Successfully"
		} 
		case "pravda": {
			doc.Find(".mdc-companies-item-title").Each(func(i int, sel *goquery.Selection) {
				// get/set title
				titleHtml     := sel.Find("span > a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
				errorService.Check(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml    := sel.Find(".mdc-companies-item-rating > span")
					rateText, _ := rateHtml.Attr("data-rating")
					setSumRate(rateText)
				}
			})

			state = "It's Parsed Successfully"
		} 
		case "spasibo": {
			doc.Find("table.items tbody tr").Each(func(i int, sel *goquery.Selection) {
				// get/set title
				titleHtml     := sel.Find("td.left > .name > a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company  + "\"\\s\\(")
				errorService.Check(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml      := sel.Find("div.stars")
					rateText, _   := rateHtml.Attr("data-fill")
					setSumRate(rateText)
					// get/set reviews
					reviewsHtml   := sel.Find("td.num")
					numReviews	   = setSumReviews(reviewsHtml.Text())
				}
			})

			state = "It's Parsed Successfully"
		} 
		case "indeedcom": {
			doc.Find(".cmp-CompanyWidget:first-child").Each(func(i int, sel *goquery.Selection) {
				// get/set title
				titleHtml     := sel.Find(".cmp-CompanyWidget-name")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile(qFeedback.Company)
				errorService.Check(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml := sel.Find(".cmp-CompanyWidget-rating-link")
					setSumRate(rateHtml.Text())
				}
			})

			state = "It's Parsed Successfully"
		}
		case "yelp": {
			doc.Find(".mainAttributes__373c0__1r0QA").Each(func(i int, sel *goquery.Selection) {
				// get/set title
				titleHtml     := sel.Find(".businessName__373c0__1fTgn h3 a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile(qFeedback.Company)
				errorService.Check(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml      := sel.Find("div.attribute__373c0__1hPI_ span div")
					rateText, _   := rateHtml.Attr("aria-label")
					rateText	      = regexp.MustCompile("[\\D]*").ReplaceAllString(rateText, "")
					setSumRate(rateText)

					// get/set reviews
					reviewsHtml   := sel.Find("div.attribute__373c0__1hPI_ .reviewCount__373c0__2r4xT")
					reviewsText   := regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
					numReviews	   = setSumReviews(reviewsText)
				}
			})

			state = "It's Parsed Successfully"
		}
		case "tripadvisorus": {
			// get All HTML
			html, err 				:= doc.Html()
			errorService.Check(err)
			html 		  				 = strings.ToLower(html)

			// get/set title
			rateViewTripExp, err := regexp.Compile(`<span>`+ qFeedback.Company +`<\/span>.*&#34;}\" data-israteable=\"true\">`)
			errorService.Check(err)
			listRateViewTrip     := rateViewTripExp.FindAllString(html, -1)

			for _, val := range listRateViewTrip {
				// get/set rate
				valExp, err	  := regexp.Compile(`\brating&#34;:&#34;[\d\.]*\b`)
				errorService.Check(err)
				numExp, err	  := regexp.Compile(`\d\.?\d?$`)
				errorService.Check(err)
				listVal  	  := valExp.FindAllString(val, -1)
				rateText 	  := numExp.FindAllString(listVal[0], -1)
				if len(rateText) != 0 {
					setSumRate(rateText[0])
				}
				
				// get/set reviews
				valExp, err    = regexp.Compile(`count&#34;:&#34;\d?\b`)
				errorService.Check(err)
				listVal        = valExp.FindAllString(val, -1)
				reviewsText   := numExp.FindAllString(listVal[0], -1)

				if len(reviewsText) != 0 {
					numReviews = setSumReviews(reviewsText[0])
				}
			}

			state = "It's Parsed Successfully"
		}
	default: {
			state = "Specified '" + title + "' Service has been not found"
			break
		}
	}

	if math.IsNaN(sumRate/float64(numRate)) {
		return 0, numReviews, state
	} else {
		return sumRate/float64(numRate), numReviews, state
	}
}

// Trim spaces ? заменить на strings.Trim!
func trim(text string) string {
	text = regexp.MustCompile("[\\s\\t\\n]*").ReplaceAllString(text, "")
	return text
}


