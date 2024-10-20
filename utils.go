package goquadac

import (
	"fmt"
	"strconv"
)

func StringtoI64(val string) int64 {
	i64, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return i64
}

func I64toString(value int64) string {
	return fmt.Sprintf("%d", value)
}

func BooltoString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func PanicOnError(message string, err error) {
	if err != nil {
		fmt.Println(message)
		panic(err)
	}
}
