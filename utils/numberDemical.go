package utils

import (
	"fmt"
	"strconv"
)

// float64Decimal2 float64取两位小数
func float64Decimal2(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}
