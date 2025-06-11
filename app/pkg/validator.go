package pkg

import (
	"go-initializr/app/common"

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
	switch common.DBType(dbType) {
	case common.DB_POSTGRES, common.DB_MYSQL, common.DB_SQLSERVER:
		return true
	}
	return false
}

func validateFrameworkType(fl valid.FieldLevel) bool {
	frameworkType := fl.Field().String()
	switch common.FrameWorkType(frameworkType) {
	case common.FRAMEWORK_ECHO, common.FRAMEWORK_GIN:
		return true
	}
	return false
}
