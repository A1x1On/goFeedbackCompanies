package feedBackAux

import (
	//"fmt"
	"github.com/PuerkitoBio/goquery"
	"gov/backend/models"
	"gov/backend/common/helpers"
	"strings"
	"math"
	"regexp"
	"strconv"
)

func ParseService(doc *goquery.Document, qFeedback *models.FeedbackQueryModel, title string) (float64, int, string) {
	state		     := "Unknown Status"
	sumRate		  := 0.0
	numRate       := 0
	numReviews    := 0
	selKey 		  := 0
	setSumRate    := func(text string){
		parsedRate, err := strconv.ParseFloat(trim(text), 64)
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

	// take one
	switch title {
		case "flamp": {
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
		case "yell": {
			doc.Find("div.companies__item-content").Each(func(i int, sel *goquery.Selection) {	
				selKey = i + 1
				// get/set title
				titleHtml     := sel.Find(".companies__item-title-text")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml 	  := sel.Find("span.rating__value")
					rateText 	  := trim(strings.ToLower(rateHtml.Text()))
					setSumRate(rateText)
					// get/set reviews
					reviewsHtml   := sel.Find("span.rating__reviews > span")
					numReviews	   = getSumReviews(reviewsHtml.Text())
				}
			})

			state = getState(selKey)
		} 
		case "apoi": {
			doc.Find("div.w_star").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/set title
				titleHtml     := sel.Find(".m_title > .flw a")
				titleText     := strings.ToLower(titleHtml.Text())
				titleExp, err := regexp.Compile("\\B\\s?" + qFeedback.Company + "\\s?\\B")
				helpers.CheckError(err)

				if titleExp.MatchString(titleText) {
					// get/set rate
					rateHtml      := sel.Find(".img_p > .p_f_s")
					rateExp, err  := regexp.Compile("[\\d\\.]*")
					rateText      := rateExp.FindAllString(trim(rateHtml.Text()), -1)
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
		case "pravda": {
			doc.Find(".mdc-companies-item-title").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/set title
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
		case "spasibo": {
			doc.Find("table.items tbody tr").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/set title
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
		case "indeedcom": {
			doc.Find(".cmp-CompanyWidget:first-child").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/set title
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
		case "yelp": {
			doc.Find(".mainAttributes__373c0__1r0QA").Each(func(i int, sel *goquery.Selection) {
				selKey = i + 1
				// get/set title
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
		case "tripadvisorus": {
			// get All HTML
			html, err 				:= doc.Html()
			helpers.CheckError(err)
			html 		  				 = strings.ToLower(html)

			// get/set title
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
		// case "yellowpages": {
		// 	// get All HTML
		// 	html, err 		:= doc.Html()
		// 	errorHelper.Check(err)

		// 	//<div class="result-rating one  "><span class="count">(2)</span>
	
		// 	html 		  		= strings.ToLower(html)

		// 	// get/set rate
		// 	rateText      := ""
		// 	valExp, err	  := regexp.Compile(`<div class="result-rating\s\w*\b`)
		// 	errorHelper.Check(err)
		// 	listRate  	  := valExp.FindAllString(html, -1)

		// 	fmt.Println("---------listRate------------", listRate)

		// 	if len(listRate) != 0 {
		// 		rateText = regexp.MustCompile(`<div class="result-rating\\s`).ReplaceAllString(listRate[0], "")
		// 		sNum	  := ""

		// 		if rateText == "one" {
		// 			sNum = "1"
		// 		} else if rateText == "two" {
		// 			sNum = "2"
		// 		} else if rateText == "three" {
		// 			sNum = "3"
		// 		} else if rateText == "four" {
		// 			sNum = "4"
		// 		} else if rateText == "five" {
		// 			sNum = "5"
		// 		}

		// 		setSumRate(sNum)
		// 	}

		// 	// get/set reviews
		// 	reviewsText   := ""
		// 	numExp, err	  := regexp.Compile(`\d\.?\d?\)`)
		// 	errorHelper.Check(err)
		// 	listReviews   := numExp.FindAllString(html, -1)
		// 	if len(listReviews) != 0 {
		// 		reviewsText    = regexp.MustCompile(`[\\D]*`).ReplaceAllString(listReviews[0], "")
		// 		numReviews	   = getSumReviews(reviewsText)
		// 	}

		// 	fmt.Println("-----------rateText----------", rateText)
		// 	fmt.Println("----------reviewsText-----------", reviewsText)

		// 	state = "It's Parsed Successfully"
		// }
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

func trim(text string) string {
	text = regexp.MustCompile("[\\s\\t\\n]*").ReplaceAllString(text, "")
	return text
}