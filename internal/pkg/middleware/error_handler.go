package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	"github.com/toji-dev/go-shop/internal/pkg/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var appErr *apperror.AppError
			if errors.As(err, &appErr) {
				switch appErr.Type {
				case apperror.TypeNotFound:
					response.NotFound(c, string(appErr.Code), appErr.Message)
				case apperror.TypeValidation:
					response.BadRequest(c, string(appErr.Code), appErr.Message, appErr.Unwrap().Error())
				case apperror.TypeConflict:
					response.Conflict(c, string(appErr.Code), appErr.Message)
				case apperror.TypeUnauthorized:
					response.Unauthorized(c, string(appErr.Code), appErr.Message)
				case apperror.TypeForbidden:
					response.Forbidden(c, string(appErr.Code), appErr.Message)
				case apperror.TypeDependencyFailure:
					response.ServiceUnavailable(c, string(appErr.Code), appErr.Message)
				default:
					response.InternalServerError(c, string(appErr.Code), appErr.Message)
				}
			} else {
				response.InternalServerError(c, string(apperror.CodeInternal), "An unexpected error occurred")
			}
		}
	}
}
