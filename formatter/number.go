package formatter

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/params"
)

func getFirstDigitAfterDecimal(num float64) string {
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
		return getFirstDigitAfterDecimal(number)
	}
	return fmt.Sprintf("%.0f", number)
}

func FormatTokenAmount(tokenAmount string, tokenDecimal int) string {
	amount := ConvertBigIntString(tokenAmount, tokenDecimal)
	amountFloat, _ := amount.Float64()
	return FormatTokenDigit(amountFloat)
}

func ConvertBigIntString(amount string, decimal int) *big.Float {
	decimals := decimal
	if decimal == 0 {
		decimals = 18
	}
	wei := new(big.Int)
	wei.SetString(amount, 10)
	return WeiToEther(wei, int(decimals))
}

func WeiToEther(wei *big.Int, decimals ...int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)

	var e *big.Float
	if len(decimals) == 0 {
		e = big.NewFloat(params.Ether)
	} else {
		e = big.NewFloat(math.Pow(10, float64(decimals[0])))
	}
	return f.Quo(fWei.SetInt(wei), e)
}

// FloatToString convert float to big int string with given decimal
// Ignore negative float
// example: FloatToString("0.000000000000000001", 18) => "1"
func FloatToString(s string, decimal int64) string {
	c, _ := strconv.ParseFloat(s, 64)
	if c < 0 {
		return "0"
	}
	bigval := new(big.Float)
	bigval.SetFloat64(c)

	d := new(big.Float)
	d.SetInt(big.NewInt(int64(math.Pow(10, float64(decimal)))))
	bigval.Mul(bigval, d)

	r := new(big.Int)
	bigval.Int(r) // store converted number in r
	return r.String()
}

// Cmp compare x and y and returns:
//
//	-1 if x <  y
//	 0 if x == y
//	+1 if x >  y
func Cmp(x, y string) (int, error) {
	n1, ok1 := new(big.Int).SetString(x, 10)
	n2, ok2 := new(big.Int).SetString(y, 10)
	if !ok1 || !ok2 {
		return 0, errors.New("invalid x or y")
	}
	return n1.Cmp(n2), nil
}

func FormatUsdDigit(value float64) string {
	tooSmall := math.Abs(value) <= 0.01

	if tooSmall {
		return "$0.01"
	}

	sign := ""
	if value < 0 {
		sign = "-"
	}

	fractionDigit := 0
	if value < 100 {
		fractionDigit = 2
	}
	params := FormatParam{
		Value:            strconv.FormatFloat(math.Abs(value), 'f', -1, 64),
		FractionDigits:   fractionDigit,
		ScientificFormat: true,
		Shorten:          math.Abs(value) >= 100,
	}

	num := formatDigit(params)

	return fmt.Sprintf("%s$%s", sign, num)
}

func FormatUsdPriceDigit(value float64) string {
	fractionDigit := 0
	if value < 100 {
		fractionDigit = 2
	}
	params := FormatParam{
		Value:            strconv.FormatFloat(math.Abs(value), 'f', -1, 64),
		FractionDigits:   fractionDigit,
		ScientificFormat: true,
		TakeExtraDecimal: 1,
		Shorten:          math.Abs(value) >= 100,
	}

	sign := ""
	if value < 0 {
		sign = "-"
	}

	num := formatDigit(params)

	return fmt.Sprintf("%s$%s", sign, num)
}

func FormatTokenDigit(value float64) string {
	params := FormatParam{
		Value:            strconv.FormatFloat(math.Abs(value), 'f', -1, 64),
		FractionDigits:   2,
		Shorten:          math.Abs(value) >= 10000,
		ScientificFormat: true,
	}
	sign := ""
	if value < 0 {
		sign = "-"
	}

	num := formatDigit(params)

	return fmt.Sprintf("%s%s", sign, num)
}

func formatDigit(params FormatParam) string {
	params.TakeExtraDecimal = int(math.Max(float64(params.TakeExtraDecimal), 0))
	num := toNum(params.Value)

	// invalid number -> keeps value the same and returns
	if num == 0 {
		return params.Value
	}
	var s string
	if num > 1 {
		s = strconv.FormatFloat(num, 'f', 9, 64)
	} else {
		s = strconv.FormatFloat(num, 'f', 18, 64)
	}
	parts := strings.Split(s, ".")
	left := parts[0]
	right := ""
	if len(parts) > 1 {
		right = parts[1]
	}

	numsArr := strings.Split(right, "")
	rightStr := numsArr[0]

	extraDecimal := params.TakeExtraDecimal
	for (toNum(rightStr) == 0 || len(rightStr) < params.FractionDigits || extraDecimal > 0) && len(numsArr) > 1 {
		if toNum(rightStr) > 0 {
			extraDecimal--
			extraDecimal = int(math.Max(float64(extraDecimal), 0))
		}
		nextDigit := numsArr[1]
		numsArr = numsArr[1:]
		rightStr += nextDigit
	}

	for strings.HasSuffix(rightStr, "0") {
		rightStr = rightStr[:len(rightStr)-1]
	}

	if len(rightStr) > params.FractionDigits && strings.Count(rightStr, "0") > 0 {
		zeroes := strings.Count(rightStr, "0")
		if len(rightStr) >= zeroes+params.FractionDigits {
			rightStr = rightStr[:zeroes+params.FractionDigits]
		}
	}

	result := left

	if params.FractionDigits != 0 && len(rightStr) > 0 {
		result += "." + rightStr
	}
	if !params.Shorten || toNum(result) == 0 {
		// find digit != 0
		index := strings.IndexFunc(right, func(r rune) bool {
			return r != '0'
		})

		digits := fmt.Sprintf("0.%s", right[:index+1])
		if index >= 6 {
			digits = strconv.FormatFloat(num, 'e', -1, 64)
			return digits
		}
		return result
	}

	if params.ScientificFormat {
		return shortenScientificNotation(toNum(result))
	}

	return strconv.FormatFloat(toNum(result), 'f', 1, 64)
}

func shortenScientificNotation(number float64) string {
	// Get the power of 10 for the number
	power := int(math.Log10(number))
	base := number / math.Pow10((power/3)*3)

	// Append the appropriate suffix (e.g., k for thousands)
	var suffix string
	switch power {
	case 0, 1:
		return strconv.FormatFloat(number, 'f', 0, 64)
	case 2:
		return strconv.FormatFloat(number, 'f', -1, 64)
	case 3, 4, 5:
		suffix = "K"
	case 6, 7, 8:
		suffix = "M"
	default:
		suffix = "B"
	}

	return strconv.FormatFloat(base, 'f', 1, 64) + suffix
}

func toNum(val interface{}) float64 {
	str := fmt.Sprintf("%v", val)
	str = strings.ReplaceAll(str, ",", "")
	if num, err := strconv.ParseFloat(str, 64); err == nil {
		return num
	}
	return 0
}
