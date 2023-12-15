package handler

import (
	custErr "github.com/deck/internal/app/error"
	"github.com/gin-gonic/gin"
	"net/http"
)

func serveHttpError(ctx *gin.Context, err error) {
	var status int
	var message string
	switch err := err.(type) {
	case *custErr.Error:
		status = err.Kind()
		message = err.Message()
	default:
		status = http.StatusInternalServerError
		message = "something went wrong"
	}

	ctx.JSON(status, map[string]string{"error": message})
}
