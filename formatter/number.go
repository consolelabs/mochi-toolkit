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
	return FormatNumberDecimal(amountFloat)
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
