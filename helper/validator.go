package helper

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func IsValidInput(input string) bool {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) == 0 {
		return false
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]`).MatchString(trimmed) {
		return false
	}

	validInputPattern := `^[a-zA-Z0-9\s!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>\/?\\|~` + "`" + `]+$`
	matched, _ := regexp.MatchString(validInputPattern, trimmed)
	if !matched {
		return false
	}

	containsAlphanumericPattern := `[a-zA-Z0-9]+`
	alphanumericMatched, _ := regexp.MatchString(containsAlphanumericPattern, trimmed)
	return alphanumericMatched
}

func ValidateUsername(username string) bool {
	trimmed := strings.TrimSpace(username)
	return len(trimmed) >= 6 && IsValidInput(trimmed)
}
