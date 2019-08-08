package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseBalance(balance string) (float64, error) {
	if balance == "" {
		return 0, nil
	}
	bl_arr := strings.Split(balance, " ")
	bl,err := strconv.ParseFloat(bl_arr[0], 64)
	return bl,err
}

func FormatBalance(balance float64, symbol string) string {
	return fmt.Sprintf("%.5f", balance) + " " + symbol
}
