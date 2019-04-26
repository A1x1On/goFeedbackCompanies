package public

import (
	"gov/backend/common/helper"
	"gov/backend/interfaces"
	"github.com/pkg/errors"
	"gov/backend/services"
	"gov/backend/models"
	"strconv"
	"bufio"
	"fmt"
	"os"
)

var feedbackService interfaces.IFeedbackService = &services.FeedbackService{}

func Index() {
	scanner   := bufio.NewScanner(os.Stdin)
	console   := &models.ConsoleModel{IsQuite: false, Step: 1}

	showMsg(console, "") // display first text instruction into the console
	
	// begin keyboard listening ... 
	for scanner.Scan() {
		switch scanner.Text() {
			case "quite" : console.IsQuite = true
			case "q"	 	 : console.IsQuite = true
			default	    : execute(console, scanner.Text()) // pick the appropriate console step in the condition blocks
		}
		
		if console.IsQuite { // if IsQuite == true do console exit
			break
		} else {
			showMsg(console, scanner.Text()) // display text instruction for next step
		}
	}

	helper.IfError(scanner.Err(), "by Enter keyboard key got error (for scanner.Scan()")
}

func showMsg(console *models.ConsoleModel, text string) {
	switch console.Step {
	case 1: {
		temp := "1 - 'flampRU', 2 - 'yellRU', 3 - 'apoiMoscow', 4 - 'pravdaRU', 5 - 'spasiboRU'\n6 - 'indeedUS', 7 - 'yelpWashington', 8 - 'tripadWashington', 9 - 'bbbUS', 10 - 'yellowWashington'\n11 - 'yelpBritan', 12 - 'yelpNorway', 13 - 'yelpPoland', 14 - 'yelpSpain', 15 - 'yelpDenmark'\n16 - 'otzyvUA'"
		fmt.Println("-----------------------\n|FEEDBACK APP IS READY|\n-----------------------\nEnter Service Id, please:\nYou can enter: \n" + temp)
	}
	case 2:
		fmt.Println("Service Id '" + text)	
		fmt.Println("Enter Company, please: ")
	default:
		helper.IfError(errors.New("Unknown Step"), "Switch default triggered (showMsg))")
	}
}

func execute(console *models.ConsoleModel, textKey string) {
	if console.Step == 1 { 
		id, err := strconv.Atoi(textKey)
		helper.IfError(err, "can't (strconv.Itoa(textKey)) into [id]")
		console.ServiceId = id
		console.Step 		= 2
	} else if console.Step == 2 { // if are All available zones
		feedback, errorCode, errorMsg := feedbackService.GetReviewService(&models.FeedbcakParamsModel{
			Company   : textKey,
			ServiceId : console.ServiceId,
		})

		fmt.Println("============== Feedbacks have been prepared =================")
		fmt.Println("----------feedback-----------", feedback)
		fmt.Println("---------errorCode------------", errorCode)
		fmt.Println("---------errorMsg------------", errorMsg)
	
		console.IsQuite = true // console exit
	} else {
		console.Step 	         = 10
	}
}
