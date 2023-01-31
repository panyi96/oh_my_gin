package middlewares

import (
	"github.com/go-playground/validator/v10"
	"ohmygin/pojo"
	"regexp"
)

func ValidatePw(fl validator.FieldLevel) bool {
	if matched, _ := regexp.MatchString(`[A-Za-z\d]*`, fl.Field().String()); matched {
		return true
	}
	return false
}

func ValidateSize(fl validator.StructLevel) {
	users := fl.Current().Interface().(pojo.Users)
	if users.UserListSize != len(users.UserList) {
		fl.ReportError(users.UserListSize, "userListSize", "UserListSize", "size必须等于数组长度", "")
	}
}
