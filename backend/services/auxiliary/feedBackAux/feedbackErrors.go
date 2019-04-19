package feedBackAux

import (
	"gov/backend/common/helper"
	"gov/backend/models"
	"regexp"
	"strconv"
	"fmt"
)


type pasingErrorInstruction struct {

}

var httpCodes = map[int]string{
	100 : "Continue",
	101 : "Switching Protocols",
	102 : "Processing",

	200 : "OK",
	201 : "Created",
	202 : "Accepted",
	203 : "Non-Authoritative Information",
	204 : "No Content",
	205 : "Reset Content",
	206 : "Partial Content",
	207 : "Multi-Status",
	208 : "Already Reported",
	226 : "IM Used",

	300 : "Multiple Choices",
	301 : "Moved Permanently",
	302 : "Moved Temporarily",
	303 : "See Other",
	304 : "Not Modified",
	305 : "Use Proxy",
	306 : "Reserved",
	307 : "Temporary Redirect",
	308 : "Permanent Redirect",

	400 : "Bad Request",
	401 : "Unauthorized",
	402 : "Payment Required",
	403 : "Forbidden",
	404 : "Not Found",
	405 : "Method Not Allowed",
	406 : "Not Acceptable",
	407 : "Proxy Authentication Required",
	408 : "Request Timeout",
	409 : "Conflict",
	410 : "Gone",
	411 : "Length Required",
	412 : "Precondition Failed",
	413 : "Payload Too Large",
	414 : "URI Too Long",
	415 : "Unsupported Media Type",
	416 : "Range Not Satisfiable",
	417 : "Expectation Failed",
	418 : "I’m a teapot",
	419 : "Authentication Timeout",
	421 : "Misdirected Request",
	422 : "Unprocessable Entity",
	423 : "Locked",
	424 : "Failed Dependency",
	426 : "Upgrade Required",
	428 : "Precondition Required",
	429 : "Too Many Requests",
	431 : "Request Header Fields Too Large",
	449 : "Retry With",
	451 : "Unavailable For Legal Reasons",
	499 : "Client Closed Request",
	500 : "Internal Server Error",
	501 : "Not Implemented",
	502 : "Bad Gateway",
	503 : "Service Unavailable",
	504 : "Gateway Timeout",
	505 : "HTTP Version Not Supported",
	506 : "Variant Also Negotiates",
	507 : "Insufficient Storage",
	508 : "Loop Detected",
	509 : "Bandwidth Limit Exceeded",
	510 : "Not Extended",
	511 : "Network Authentication Required",
	520 : "Unknown Error",
	521 : "Web Server Is Down",
	522 : "Connection Timed Out",
	523 : "Origin Is Unreachable",
	524 : "A Timeout Occurred",
	525 : "SSL Handshake Failed",
	526 : "Invalid SSL Certificate",
}

var parsingCodes = map[int]string{
	1000 : "Tag has not been Found for the ",
	1001 : "Switching Protocols",
	1002 : "Processing",
}

func setHttpErrorByHtml(text string, errorState *models.ErrorStateModel){
	for key, val := range httpCodes {
		matched, err := regexp.MatchString(strconv.Itoa(key), text)
		helper.CheckError(err)

		if matched {
			*errorState = models.ErrorStateModel{Message: val, Code: key}
			break
		}
	}
}

func setHttpErrorByCode(code int, errorState *models.ErrorStateModel){
	for key, val := range httpCodes {
		if key == code {
			*errorState = models.ErrorStateModel{Message: val, Code: key}
			break
		}
	}
}

func setParsingErrorByCode(code int, tagType string, errorState *models.ErrorStateModel){
	fmt.Println("----------code-----------", code)

	for key, val := range parsingCodes {
		if key == code {
			*errorState = models.ErrorStateModel{Message: val + tagType, Code: key}
			break
		}
	}
}

// func pasingErrorAnalyzer(numFound int, tagType string, errorState *models.ErrorStateModel){
// 	// title = 

// 	fmt.Println("---------numFound------------", numFound)

// 	if numFound == 0 {

// 		switch tagType {
// 			case "title" : {
// 				setParsingError(1000, tagType, errorState)
// 			}
// 		}

// 	}
// }
