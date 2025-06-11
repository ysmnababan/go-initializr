package pkg

import (
	"go-initializr/app/initializer"

	valid "github.com/go-playground/validator/v10"
)

type customValidator struct {
	validator *valid.Validate
}

func NewCustomValidator() *customValidator {
	newValidator := valid.New()
	_ = newValidator.RegisterValidation("db-type", validateDbType)
	_ = newValidator.RegisterValidation("framework-type", validateFrameworkType)

	return &customValidator{
		validator: newValidator,
	}
}

func (cv *customValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func validateDbType(fl valid.FieldLevel) bool {
	dbType := fl.Field().String()
	switch initializer.DBType(dbType) {
	case initializer.DB_POSTGRES, initializer.DB_MYSQL, initializer.DB_SQLSERVER:
		return true
	}
	return false
}

func validateFrameworkType(fl valid.FieldLevel) bool {
	frameworkType := fl.Field().String()
	switch initializer.FrameWorkType(frameworkType) {
	case initializer.FRAMEWORK_ECHO, initializer.FRAMEWORK_GIN:
		return true
	}
	return false
}
