package formatter

import (
	"fmt"
	"strconv"
	"strings"
)

func GetFirstDigitAfterDecimal(num float64) string {
	str := strconv.FormatFloat(num, 'f', -1, 64)

	// Split the string into two parts: the integer part and the decimal part
	parts := strings.Split(str, ".")
	if len(parts) != 2 {
		return "0"
	}

	// Retrieve the decimal part
	decimalStr := parts[1]

	// find digit != 0
	index := strings.IndexFunc(decimalStr, func(r rune) bool {
		return r != '0'
	})

	digits := fmt.Sprintf("0.%s", decimalStr[:index+1])
	if index >= 6 {
		digits = strconv.FormatFloat(num, 'E', -1, 64)
	}

	return strings.ToLower(digits)
}

func FormatNumberDecimal(number float64) string {
	if number < 1000 && number >= 1 {
		return fmt.Sprintf("%.2f", number)
	}
	if number < 1 {
		return GetFirstDigitAfterDecimal(number)
	}
	return fmt.Sprintf("%.0f", number)
}

func FormatTokenAmount(tokenAmount string, tokenDecimal int) string {
	amount := ConvertBigIntString(tokenAmount, tokenDecimal)
	amountFloat, _ := amount.Float64()
	return FormatNumberDecimal(amountFloat)
}
