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
				"old_inputs": inputs,
			})
}

func (ac *AuthController) LoginPage(c fiber.Ctx) error {
	return c.Render("login", fiber.Map{
				"Title": "Login Form",
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) ForgotPasswordPage(c fiber.Ctx) error {
	return c.Render("forgot-password", fiber.Map{
				"Title": "Forgot Password Form",
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) ResetPasswordPage(c fiber.Ctx) error {
	return c.Render("reset-password", fiber.Map{
				"Title": "Reset Password Form",
				"csrf": csrf.TokenFromContext(c),
				"token": "",
			})
}

func (ac *AuthController) PasswordConfirmationPage(c fiber.Ctx) error {
	return c.Render("password-confirmation", fiber.Map{
				"Title": "Password Confirmation Form",
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) TwoFactorAuthChallengePage(c fiber.Ctx) error {
	return c.Render("2fa-challenge", fiber.Map{
				"Title": "2FA Challenge Form",
				"csrf": csrf.TokenFromContext(c),
			})
}

func (ac *AuthController) TwoFactorAuthRecoveryCodesPage(c fiber.Ctx) error {
	return c.Render("2fa-recovery-codes", fiber.Map{
				"Title": "2FA Recovery Codes Form",
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
			return c.Redirect().Status(fiber.StatusSeeOther).With( "error", "invalid form submission" ).Back()
		}
		spew.Dump( user )

	// Validate user inputs
		if err := ac.validator.Struct( user ); err != nil {
			var validationErrors validator.ValidationErrors
			errors.As( err, &validationErrors )

			// Translate errors into a map: map[field_name]error_message
				translatedErrors := validationErrors.Translate(ac.translator)

			// Convert it to JSON
				validationErrorsInJSON, _ := json.Marshal( translatedErrors )
				return c.Redirect().Status(fiber.StatusSeeOther).WithInput().With( "validation_errors", string(validationErrorsInJSON) ).Back()
		}

		


	// Call actions
		// if err := ac.auth.CreateUser(user); err != nil {
		// 	// handle ErrDuplicateEmail etc.
		// }
		// if err := ac.auth.SendEmailVerificationNotification(user); err != nil {
		// 	// handle mail failure
		// }

	// 3. respond: redirect to "check your email" page, or JSON message


	//----------------------------------------
	return c.SendString( "Process Register" )
}

func (ac *AuthController) Login(c fiber.Ctx) error {
	// return c.Redirect().To( "/auth/login" )
	return c.SendString( "Process Login" )
}

func (ac *AuthController) ForgotPassword(c fiber.Ctx) error {
	// return c.Redirect().To( "/auth/login" )
	return c.SendString( "Process ForgotPassword" )
}