package form

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Form struct {
	errors map[string][]string
}

// BindAndValidate binds the form to the request and validates it.
func (f *Form) BindAndValidate(ctx echo.Context, form any) error {

	if err := ctx.Bind(form); err != nil {
		return err
	}

	if err := ctx.Validate(form); err != nil {
		f.setErrorMessage(err)
		return err
	}
	return nil
}

func (f *Form) AddError(field string, message string) {
	if f.errors == nil {
		f.errors = make(map[string][]string)
	}

	f.errors[field] = append(f.errors[field], message)
}

func (f *Form) Errors() map[string][]string {
	return f.errors
}

func (f *Form) IsValid() bool {
	if f.errors == nil {
		return true
	}
	return len(f.errors) == 0
}

func (f *Form) GetErrors(field string) []string {
	if f.errors == nil {
		return nil
	}
	return f.errors[field]
}

func (f *Form) setErrorMessage(err error) {

	ves, ok := err.(validator.ValidationErrors)
	if !ok {
		return
	}

	for _, ve := range ves {
		var message string

		switch ve.Tag() {
		case "required":
			message = "required"
		case "email":
			message = "invalid email"
		case "min":
			message = fmt.Sprintf("min %s characters long", ve.Param())
		case "max":
			message = fmt.Sprintf("max %s characters long", ve.Param())
		default:
			message = "invalid"
		}

		f.AddError(ve.Field(), message)

	}
}
