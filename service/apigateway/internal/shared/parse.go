package shared

import (
	"strconv"
	"strings"

	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	cardapierrors "github.com/MamangRust/monolith-payment-gateway-shared/errors/card_errors/api"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// parseQueryInt is a helper function for parsing an int query parameter.
//
// It takes a Context, key, and default value. It will return the default value
// if the query parameter is not set, or if it cannot be parsed into an int.
func ParseQueryInt(c echo.Context, key string, defaultValue int) int {
	param := c.QueryParam(key)
	if param == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(param)
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}

// parseQueryCard parses a query parameter string into a card number.
//
// It takes a Context and Logger. It will return an empty string and an error if
// the query parameter is not set or if the query parameter is empty.
func ParseQueryCard(c echo.Context, logger logger.LoggerInterface) (string, error) {
	cardNumber := strings.TrimSpace(c.QueryParam("card_number"))
	if cardNumber == "" {
		logger.Error("card number is empty")
		return "", cardapierrors.ErrApiInvalidCardNumber(c)
	}
	return cardNumber, nil
}

// parseQueryMonth is a helper function for parsing a query parameter string into a month.
//
// It takes a Context and Logger. It will return an error if the query parameter is not set or
// if the query parameter is not a valid month (1-12).
func ParseQueryMonth(c echo.Context, logger logger.LoggerInterface) (int, error) {
	monthStr := c.QueryParam("month")
	month, err := strconv.Atoi(monthStr)

	if err != nil || month < 1 || month > 12 {
		logger.Error("invalid month", zap.String("month", monthStr))
		return 0, cardapierrors.ErrApiInvalidMonth(c)
	}

	return month, nil
}

// parseQueryYear parses a query parameter string into a year.
//
// It takes a Context and Logger. It will return an error if the query parameter
// is not set, cannot be parsed into an integer, or if the year is less than 2023.
func ParseQueryYear(c echo.Context, logger logger.LoggerInterface) (int, error) {
	yearStr := c.QueryParam("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2023 {
		logger.Error("invalid year", zap.String("year", yearStr))
		return 0, cardapierrors.ErrApiInvalidYear(c)
	}
	return year, nil
}
