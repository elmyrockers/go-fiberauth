package controller

import (
	"errors"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	fiberauth "github.com/elmyrockers/go-fiberauth/session"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	
	"github.com/elmyrockers/go-fiberauth/examples/session/form"
	"github.com/davecgh/go-spew/spew"
)

type AuthController struct {
	auth *fiberauth.Auth
	validator *validator.Validate
	translator ut.Translator
}

func NewAuth(auth *fiberauth.Auth) *AuthController {
	validator := validator.New()

	enLocale := en.New()
	uni := ut.New( enLocale )
	translator, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validator, translator)

	return &AuthController{ auth: auth, validator: validator, translator: translator }
}

// HTTP METHOD: GET
func (ac *AuthController) RegisterPage(c fiber.Ctx) error {
	var fieldErrors map[string]string

	// Fetch JSON string from flash message
		error := c.Redirect().Message("error").Value
		validationErrors := c.Redirect().Message("validation_errors").Value
		if validationErrors != "" {
			_ = json.Unmarshal([]byte(validationErrors), &fieldErrors)
		}

	// Remap old inputs
		oldInputs := c.Redirect().OldInputs()
		inputs := make(map[string]string, len(oldInputs))
		for _, item := range oldInputs {
			inputs[ item.Key ] = item.Value
		}

	return c.Render("register", fiber.Map{
				"csrf": csrf.TokenFromContext(c),
				"validation_errors": fieldErrors,
				"error":error,
				"old_inputs": inputs,
			})
}

func (ac *AuthController) LoginPage(c fiber.Ctx) error {
	success := c.Redirect().Message( "success" ).Value
	return c.Render("login", fiber.Map{
				"success": success,
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) ForgotPasswordPage(c fiber.Ctx) error {
	return c.Render("forgot-password", fiber.Map{
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) ResetPasswordPage(c fiber.Ctx) error {
	return c.Render("reset-password", fiber.Map{
				"csrf": csrf.TokenFromContext(c),
				"token": "",
			})
}

func (ac *AuthController) PasswordConfirmationPage(c fiber.Ctx) error {
	return c.Render("password-confirmation", fiber.Map{
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) TwoFactorAuthChallengePage(c fiber.Ctx) error {
	return c.Render("2fa-challenge", fiber.Map{
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) TwoFactorAuthRecoveryCodesPage(c fiber.Ctx) error {
	return c.Render("2fa-recovery-codes", fiber.Map{
				"csrf": csrf.TokenFromContext(c),
				"codes": []string{
					"A1BC-DEFG",
					"H2IJ-KLMN",
					"OPQR-STUV",
					"WXYZ-1234",
					"5678-90AB",
					"CDEF-GHIJ",
					"KLMN-OPQR",
					"STUV-WXYZ",
				},
			})
}

// HTTP METHOD: POST, PUT
func (ac *AuthController) Register(c fiber.Ctx) error {
	// Parse + validate input (name, email, password, password_confirmation)
		var user form.RegisterForm
		if err := c.Bind().Body(&user); err != nil {
			return c.Redirect().With( "error", "Invalid form submission" ).Back()
		}
		spew.Dump( user )

	// Validate user inputs
		if err := ac.validator.Struct( user ); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As( err, &validationErrors )

			// Translate errors into a map: map[field_name]error_message
				translatedErrors := validationErrors.Translate(ac.translator)

			// Convert it to JSON then redirect
				validationErrorsInJSON, _ := json.Marshal( translatedErrors )
				return c.Redirect().WithInput().With( "validation_errors", string(validationErrorsInJSON) ).Back()
		}

	// Create new user
		userID, err := ac.auth.CreateUser( user.Name, user.Email, user.Password )
		if err != nil {
			var errMessage string
			if errors.Is(err, fiberauth.ErrEmailTaken) {
				errMessage = "Email is already registered"
			} else if errors.Is(err, fiberauth.ErrPasswordTooLong) {
				errMessage = "Password is too long"
			} else {
				errMessage = "Failed to create account"
			}
			return c.Redirect().WithInput().With( "error", errMessage ).Back()
		}

	// Send a notification for email verification
		if err := ac.auth.SendEmailVerificationNotification(userID,user.Name,user.Email, "http://localhost:3000/auth/verify-email"); err != nil {
			return c.Redirect().WithInput().With( "error", "Account created, but failed to send verification email" ).Back()
		}

	// Respond: redirect to login page with a success message
		return c.Redirect().With( "success", "Account created - check your email to verify before logging in" ).To( "/auth/login" )
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	return c.SendString( "Process Login" )
}

func (ac *AuthController) ForgotPassword(c fiber.Ctx) error {
	return c.SendString( "Process ForgotPassword" )
}