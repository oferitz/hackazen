package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/oferitz/hackazen/internal/validate"
	"net/http"
)

func (app *application) errorResponse(status int, message string) error {
	return fiber.NewError(status, message)
}

func (app *application) serverErrorResponse(c *fiber.Ctx, err error) error {
	app.logger.Error(err.Error())
	msg := "the server encountered a problem and could not process your request"
	return app.errorResponse(fiber.StatusInternalServerError, msg)
}

func (app *application) notFoundResponse() error {
	msg := "the requested resource could not be found"
	return app.errorResponse(fiber.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(c *fiber.Ctx) error {
	msg := fmt.Sprintf("the %s method is not supported for this resource", c.Method())
	return app.errorResponse(fiber.StatusMethodNotAllowed, msg)
}

func (app *application) badRequestResponse(err error) error {
	return app.errorResponse(fiber.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(c *fiber.Ctx, errors []*validate.ErrorResponse) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
}

func (app *application) editConflictResponse() error {
	message := "unable to update the record due to an edit conflict, please try again"
	return app.errorResponse(fiber.StatusConflict, message)
}

func (app *application) entityAlreadyExists(message string) error {
	return app.errorResponse(fiber.StatusConflict, message)
}

func (app *application) rateLimitExceededResponse() error {
	message := "rate limit exceeded"
	return app.errorResponse(fiber.StatusTooManyRequests, message)
}

func (app *application) invalidCredentialsResponse() error {
	message := "invalid authentication credentials"
	return app.errorResponse(fiber.StatusUnauthorized, message)
}

func (app *application) invalidAuthenticationTokenResponse(c *fiber.Ctx) error {
	c.Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	return app.errorResponse(fiber.StatusUnauthorized, message)
}

func (app *application) authenticationRequiredResponse() error {
	message := "you must be authenticated to access this resource"
	return app.errorResponse(http.StatusUnauthorized, message)
}

func (app *application) inactiveAccountResponse() error {
	message := "your user account must be activated to access this resource"
	return app.errorResponse(fiber.StatusForbidden, message)
}

func (app *application) notPermittedResponse() error {
	message := "your user account doesn't have the necessary permissions to access this resource"
	return app.errorResponse(fiber.StatusForbidden, message)
}
