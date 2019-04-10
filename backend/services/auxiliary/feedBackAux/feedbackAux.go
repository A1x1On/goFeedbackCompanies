package feedBackAux

import (
	"github.com/PuerkitoBio/goquery"
	"gov/backend/models"
	"gov/backend/common/helpers"
	"encoding/json"
	"strings"
	"math"
	"math/big"
	"regexp"
	"strconv"
)

type bbbRating struct {
	Search struct {
		Results []struct {
			BusinessName        string `json:"businessName"`
			Rating              string `json:"rating"`
		} `json:"results"`
	} `json:"search"`
}


func ParseService(doc *goquery.Document, qFeedback *models.FeedbackQueryModel, title string) (float64, int, string) {
	state		     := "Unknown Status"
	sumRate		  := 0.0
	numRate       := 0
	numReviews    := 0
	selKey 		  := 0
	setSumRate    := func(text string){
		parsedRate, err := strconv.ParseFloat(trimAll(text), 64)
		helpers.CheckError(err)
		if parsedRate != 0 {
			sumRate    = sumRate + parsedRate
			numRate += 1
		}
	}
	getSumReviews := func(reviewsText string) int{
		reviewsInt, _ := strconv.Atoi(reviewsText)
		return numReviews + reviewsInt
	}
	getState		  := func(iKey int) string{
		if iKey == 0 {
			html, err 		:= doc.Html()
			helpers.CheckError(err)
			code, msg := helpers.CheckIsHttpError(html)
			return `"` + code + `" ` + msg
		} else {
			return "It's Parsed Successfully"
		}
	}
	caseForYelp		:= func(){
		doc.Find(".mainAttributes__373c0__1r0QA").Each(func(i int, sel *goquery.Selection) {
			selKey = i + 1
			// get/check title
			titleHtml     := sel.Find(".heading--h3__373c0__1n4Of > a")
			titleText     := strings.ToLower(titleHtml.Text())
			titleExp, err := regexp.Compile(qFeedback.Company)
			helpers.CheckError(err)

			if titleExp.MatchString(titleText) {
				// get/set rate
				rateHtml 	  := sel.Find(".i-stars__373c0__30xVZ")
				rateText, _   := rateHtml.Attr("aria-label")
				rateExp, err  := regexp.Compile(`\d\.?\d?`)
				helpers.CheckError(err)
				listRate  	  := rateExp.FindAllString(rateText, -1)
				if len(listRate) != 0 {
					setSumRate(listRate[0])
				}
			
				// get/set reviews
				reviewsHtml   := sel.Find(".reviewCount__373c0__2r4xT")
				reviewsText	  := regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
				numReviews	   = getSumReviews(reviewsText)
			}
		})

		state = getState(selKey)
	}	

	// prepare one feedback from some found company results in one service 
	switch title {
		case "flampRU": {
			doc.Find("cat-brand-filial-rating").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/set rate
				rateText, _          := sel.Attr("rating")
				setSumRate(rateText)
				// get/set reviews
				reviewsText, _       := sel.Attr("reviews-count")
				numReviews	   		 = getSumReviews(reviewsText)
			})

			state = getState(selKey)
		} 
		case "yellRU": {
			doc.Find("div.companies__item-content").Each(func(i int, sel *goquery.Selection) {	
				selKey = i + 1
				// get/check title
				titleHtml     := sel.Find(".companies__item-title-text")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml 	  := sel.Find("span.rating__value")
					rateText 	  := trimAll(strings.ToLower(rateHtml.Text()))
					setSumRate(rateText)
					// get/set reviews
					reviewsHtml   := sel.Find("span.rating__reviews > span")
					numReviews	   = getSumReviews(reviewsHtml.Text())
				}
			})

			state = getState(selKey)
		} 
		case "apoiMoscow": {
			doc.Find("div.w_star").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/check title
				titleHtml     := sel.Find(".m_title > .flw a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml      := sel.Find(".img_p > .p_f_s")
					rateExp, err  := regexp.Compile("[\\d\\.]*")
					rateText      := rateExp.FindAllString(trimAll(rateHtml.Text()), -1)
					helpers.CheckError(err)
					setSumRate(rateText[0])
					// get/set reviews
					reviewsHtml   := sel.Find(".img_p > .numReviews")
					reviewsText	  := regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
					numReviews	   = getSumReviews(reviewsText)
				}
			})

			state = getState(selKey)
		} 
		case "pravdaRU": {
			doc.Find(".mdc-companies-item-title").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/check title
				titleHtml     := sel.Find("span > a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml    := sel.Find(".mdc-companies-item-rating > span")
					rateText, _ := rateHtml.Attr("data-rating")
					setSumRate(rateText)
				}
			})

			state = getState(selKey)
		} 
		case "spasiboRU": {
			doc.Find("table.items tbody tr").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/check title
				titleHtml     := sel.Find("td.left > .name > a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company  + "\"\\s\\(")
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml      := sel.Find("div.stars")
					rateText, _   := rateHtml.Attr("data-fill")
					setSumRate(rateText)
					// get/set reviews
					reviewsHtml   := sel.Find("td.num")
					numReviews	   = getSumReviews(reviewsHtml.Text())
				}
			})

			state = getState(selKey)
		} 
		case "indeedUS": {
			doc.Find(".cmp-CompanyWidget:first-child").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/check title
				titleHtml     := sel.Find(".cmp-CompanyWidget-name")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile(qFeedback.Company)
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml := sel.Find(".cmp-CompanyWidget-rating-link")
					setSumRate(rateHtml.Text())
				}
			})

			state = getState(selKey)
		}
		case "yelpWashington": {
			doc.Find(".mainAttributes__373c0__1r0QA").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/check title
				titleHtml     := sel.Find(".businessName__373c0__1fTgn h3 a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile(qFeedback.Company)
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml      := sel.Find("div.attribute__373c0__1hPI_ span div")
					rateText, _   := rateHtml.Attr("aria-label")
					rateText	      = regexp.MustCompile("[\\D]*").ReplaceAllString(rateText, "")
					setSumRate(rateText)

					// get/set reviews
					reviewsHtml   := sel.Find("div.attribute__373c0__1hPI_ .reviewCount__373c0__2r4xT")
					reviewsText   := regexp.MustCompile("[\\D]*").ReplaceAllString(reviewsHtml.Text(), "")
					numReviews	   = getSumReviews(reviewsText)
				}
			})

			state = getState(selKey)
		}
		case "tripadWashington": {
			// get All HTML
			html, err 				:= doc.Html()
			helpers.CheckError(err)
			html 		  				 = strings.ToLower(html)

			// get/check title
			rateViewTripExp, err := regexp.Compile(`<span>`+ qFeedback.Company +`<\/span>.*&#34;}\" data-israteable=\"true\">`)
			helpers.CheckError(err)
			listRateViewTrip     := rateViewTripExp.FindAllString(html, -1)

			for key, val := range listRateViewTrip {
				selKey = key + 1
				// get/set rate
				valExp, err	  := regexp.Compile(`\brating&#34;:&#34;[\d\.]*\b`)
				helpers.CheckError(err)
				numExp, err	  := regexp.Compile(`\d\.?\d?$`)
				helpers.CheckError(err)
				listVal  	  := valExp.FindAllString(val, -1)
				rateText 	  := numExp.FindAllString(listVal[0], -1)
				if len(rateText) != 0 {
					setSumRate(rateText[0])
				}
				
				// get/set reviews
				valExp, err    = regexp.Compile(`count&#34;:&#34;\d?\b`)
				helpers.CheckError(err)
				listVal        = valExp.FindAllString(val, -1)
				reviewsText   := numExp.FindAllString(listVal[0], -1)

				if len(reviewsText) != 0 {
					numReviews = getSumReviews(reviewsText[0])
				}
			}

			state = getState(selKey)
		}
		case "yellowWashington": {
			// get HTML
			html, err 	  := doc.Html()
			helpers.CheckError(err)
			html 		  		= strings.ToLower(html)

			// get rate & reviews as DOM string array by found companies
			rateText      := ""
			valExp, err	  := regexp.Compile(`<div class="result-rating\s\D*\d*\)`)
			helpers.CheckError(err)
			listValue  	  := valExp.FindAllString(html, -1)

			if len(listValue) == 0 { 
				break
			}

			for i, val := range listValue {
				selKey = i + 1

				// get/set rate
				rateExp, err  := regexp.Compile(`one|two|three|four|five`)
				helpers.CheckError(err)
				rateNumText   := rateExp.FindAllString(val, -1)

				if len(rateNumText) != 0 {
					if rateNumText[0] == "one" {
						rateText = "1.0"
					} else if rateNumText[0]  == "two" {
						rateText = "2.0"
					} else if rateNumText[0]  == "three" {
						rateText = "3.0"
					} else if rateNumText[0]  == "four" {
						rateText = "4.0"
					} else if rateNumText[0]  == "five" {
						rateText = "5.0"
					} else {
						rateText = "0.0"
					}
					setSumRate(rateText)
				}

				// get/set reviews
				numExp, err	  := regexp.Compile(`\d*\)`)
				helpers.CheckError(err)
				reviewsText   := numExp.FindAllString(val, -1)
				if len(reviewsText) != 0 {
					numReviews = getSumReviews(reviewsText[0][:1])
				}
			}

			state = getState(selKey)
		}
		case "bbbUS": {
			// get HTML
			html, err 	  := doc.Html()
			helpers.CheckError(err)
			html 		  		= strings.ToLower(html)

			// get title & rate as json by found companies
			sRating  	  := bbbRating{}
			rateText      := ""
			valExp, err	  := regexp.Compile(`bbbDtmData.*\}`)
			helpers.CheckError(err)
			listValue  	  := valExp.FindAllString(doc.Text(), -1)

			if len(listValue) == 0 {
				break
			}

			jsonExp, err := regexp.Compile(`\{.*`)
			helpers.CheckError(err)
			ratingJson   := jsonExp.FindAllString(listValue[0], -1)

			if len(ratingJson) == 0 {
				break
			}

			json.Unmarshal([]byte(ratingJson[0]), &sRating)

			// got json result from service document
			for i, val := range sRating.Search.Results {
				selKey			  = i + 1
				val.BusinessName = strings.ToLower(val.BusinessName)
				grade				 := string([]byte(val.Rating)[0])
				
				// get/check title
				titleText       := strings.ToLower(val.BusinessName)
				titleExp, err   := regexp.Compile("^" + qFeedback.Company + "$")
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					if grade == "F" {
						rateText = "1.0"
					} else if grade  == "D" {
						rateText = "2.0"
					} else if grade  == "C" {
						rateText = "3.0"
					} else if grade  == "B" {
						rateText = "4.0"
					} else if grade  == "A" {
						rateText = "5.0"
					} else {
						rateText = "0.0"
					}
					setSumRate(rateText)	
				}
			}

			state = getState(selKey)
		} 
		case "yelpPoland" : caseForYelp()	
		case "yelpSpain"  : caseForYelp()
		case "yelpDenmark": caseForYelp()
		case "yelpBritan" : caseForYelp()
		case "yelpNorway" : caseForYelp()
		case "otzyvUA"	   : {
			doc.Find(".otzyv_box_float").Each(func(i int, sel *goquery.Selection) {
				selKey         = i + 1
				valHtml       := sel.Find(".otzyv_item_cat1")
	
				// get/set rate
				rateExp, err  := regexp.Compile(`г\s\d\.?\d?`)
				helpers.CheckError(err)
				listRate 	  := rateExp.FindAllString(valHtml.Text(), -1)

				if len(listRate) != 0 {
					rateExp, err  := regexp.Compile(`\d\.?\d?`)
					helpers.CheckError(err)
					listRate  	  := rateExp.FindAllString(trimAll(listRate[0]), -1)
					if len(listRate) != 0 {
						setSumRate(listRate[0])
					}
				}
	
				// get/set reviews
				reviewsExp, err := regexp.Compile(`\s\d*\sо`)
				helpers.CheckError(err)
				listReviews     := reviewsExp.FindAllString(valHtml.Text(), -1)
				if len(listReviews) != 0 {
					reviewsText	:= regexp.MustCompile("[\\D]*").ReplaceAllString(strings.TrimLeft(listReviews[0], " "), "")
					numReviews	 = getSumReviews(reviewsText)
				}
			})
	
			state = getState(selKey)
		}
		default: {
			state = "Specified '" + title + "' Service has been not found"
			break
		}
	}

	if math.IsNaN(sumRate/float64(numRate)) {
		return 0, numReviews, state
	} else {
		fixedRate, err	:= strconv.ParseFloat(big.NewFloat(sumRate/float64(numRate)).Text('f', 3), 64)
		helpers.CheckError(err)

		return fixedRate, numReviews, state
	}
}

// trin all spaces and new lines
func trimAll(text string) string {
	text = regexp.MustCompile("[\\s\\t\\n]*").ReplaceAllString(text, "")
	return text
}