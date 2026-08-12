package form


type RegisterForm struct {
	Name                  string `json:"name" form:"name" validate:"required,max=255"`
	Email                 string `json:"email" form:"email" validate:"required,email,max=255"`
	Password              string `json:"password" form:"password" validate:"required,min=8"`
	PasswordConfirmation  string `json:"password_confirmation" form:"password_confirmation" validate:"required,eqfield=Password"`
}