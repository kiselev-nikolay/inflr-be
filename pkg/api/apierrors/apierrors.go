package apierrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	code   int
	status string
}

func (ae *APIError) Send(g *gin.Context) {
	g.JSON(ae.code, gin.H{
		"status": ae.status,
	})
}

var (
	MissingRequiredData APIError
	AlreadyHave         APIError
	CannotCreate        APIError
	NotFound            APIError
)

func init() {
	MissingRequiredData = APIError{http.StatusBadRequest, "missing required data"}
	AlreadyHave = APIError{http.StatusBadRequest, "already have"}
	CannotCreate = APIError{http.StatusInternalServerError, "cannot create"}
	NotFound = APIError{http.StatusNotFound, "not found"}
}
