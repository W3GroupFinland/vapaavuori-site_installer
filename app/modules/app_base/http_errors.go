package app_base

import (
	"fmt"
	"log"
	"net/http"
)

type Http struct {
}

type HttpStatus int

var (
	HttpErrorMessages = map[int]string{
		http.StatusAccepted:            "202 Ok.",
		http.StatusUnauthorized:        "401 Unauthorized.",
		http.StatusForbidden:           "403 Forbidden.",
		http.StatusNotFound:            "404 Not found.",
		http.StatusInternalServerError: "500 Internal server error.",
	}
)

func (h *Http) Error(httpStatus int, errorMsg ...string) string {
	return HttpErrorMsg(httpStatus, errorMsg...)
}

func HttpErrorMsg(httpStatus int, errorMsgs ...string) string {
	var (
		msg    string
		exists bool
		errMsg string
	)

	var appMsg string

	for i := 0; len(errorMsgs) > i; i++ {
		appMsg += "\n"
		appMsg += errorMsgs[i]
	}

	if msg, exists = HttpErrorMessages[httpStatus]; !exists {
		errMsg = fmt.Sprintf("Error code %v%v", httpStatus, appMsg)
		log.Printf("No error message found for status code %v\n", httpStatus)
	} else {
		errMsg = fmt.Sprintf("%v%v", msg, appMsg)
	}

	log.Println(errMsg)

	return errMsg
}
