package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"woojiahao.com/hermes/internal"
)

func ginError(ctx *gin.Context, errorCode int, message any) {
	ctx.JSON(errorCode, errorBody{errorCode, message})
}

func internalSeverError(ctx *gin.Context) {
	ginError(ctx, http.StatusInternalServerError, "Internal server error")
}

func notFound(ctx *gin.Context, message string) {
	ginError(ctx, http.StatusNotFound, message)
}

func badRequestValidation(ctx *gin.Context, bindingError error) {
	var ve validator.ValidationErrors
	if errors.As(bindingError, &ve) {
		out := internal.Map(ve, func(field validator.FieldError) errorField {
			message := ""
			switch field.Tag() {
			case "required":
				message = "This field is required"
			case "min":
				message = "This field has a minimum necessary length/size"
			case "email":
				message = "This field must be an email"
			default:
				message = field.Tag()
			}
			return errorField{field.Field(), message}
		})

		ginError(ctx, http.StatusBadRequest, out)
	}
}

func badRequest(ctx *gin.Context, message string) {
	ginError(ctx, http.StatusBadRequest, message)
}
