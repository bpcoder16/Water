package utils

import "fmt"

func MaskedName(name string) string {
	maskedStr := "***"
	if len(name) <= 1 {
		return name + maskedStr
	} else {
		return fmt.Sprintf("%c%s%c", name[0], maskedStr, name[len(name)-1])
	}
}
