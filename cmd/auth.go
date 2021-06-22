package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oferitz/hackazen/internal/data"
	"github.com/oferitz/hackazen/internal/validate"
	"time"
)

func (app *application) signupHandler(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email" validate:"required,email,min=6,max=32"`
		Password string `json:"password" validate:"required,min=8,max=32"`
	}

	if err := c.BodyParser(&input); err != nil {
		return app.serverErrorResponse(c, err)
	}

	validationErrors := validate.ValidateStruct(input)
	if validationErrors != nil {
		return app.failedValidationResponse(c, validationErrors)
	}

	user := &data.User{
		Email:     input.Email,
		Activated: false,
	}

	err := user.Password.Set(input.Password)
	if err != nil {
		return app.serverErrorResponse(c, err)
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			return app.entityAlreadyExists("Someone’s already using that email.")
		default:
			return app.serverErrorResponse(c, err)
		}
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		return app.serverErrorResponse(c, err)
	}

	// send welcome email in the background
	app.background(func() {
		tmplData := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.gohtml", tmplData)
		if err != nil {
			app.logger.Error(err, nil)
		}
	})

	return c.Send([]byte("yofi"))

}
