package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message    string
	Code       int
	StatusCode int
}

func SendErrorResponse(c *gin.Context, e ErrorResponse) {
	c.JSON(e.StatusCode, gin.H{
		"message": e.Message,
		"error":   e.Code,
	})
	c.Abort()
}

var NOT_FOUND_ERROR = ErrorResponse{"Not Found", http.StatusNotFound, http.StatusNotFound}
var INVALID_REQUEST_METHOD_ERROR = ErrorResponse{"Invalid request method", http.StatusMethodNotAllowed, http.StatusMethodNotAllowed}
var INTERNAL_SERVER_ERROR = ErrorResponse{"Internal server error", http.StatusInternalServerError, http.StatusInternalServerError}
var TO_MANY_REQUESTS = ErrorResponse{"Too many request, please try again later", http.StatusTooManyRequests, http.StatusTooManyRequests}

// 10xx error code for authentication related
var INVALID_REGISTER_FORM = ErrorResponse{"Invalid register form", 1010, http.StatusBadRequest}
var USER_EXISTS_FOR_SIGNUP = ErrorResponse{"User exists login please", 1011, http.StatusBadRequest}
var INVALID_EMAIL_ADDRESS = ErrorResponse{"Invalid email address", 1012, http.StatusBadRequest}
var INVALID_USERNAME = ErrorResponse{"Invalid username", 1013, http.StatusBadRequest}
var FAILED_TO_CREATE_USER = ErrorResponse{"Failed to create user", 1014, http.StatusInternalServerError}

var INVALID_CREDENTIALS = ErrorResponse{"Invalid login credentials", 1021, http.StatusBadRequest}

var NO_COOKIE_FOUND = ErrorResponse{"No cookie found", 1030, http.StatusForbidden}
var COOKIE_EXPIRED = ErrorResponse{"Invalid cookie or cookie expired", 1031, http.StatusForbidden}


// errors for backend to handle
type Error string

func (e Error) Error() string {
	return string(e)
}

var UnexpectedSigningMethodErr = Error("Unexpected signing method")
