package utils

func MaskedName(name string) string {
	maskedStr := "***"
	runeNames := []rune(name)
	if len(runeNames) <= 1 {
		return name + maskedStr
	} else {
		return string(runeNames[0]) + maskedStr + string(runeNames[len(runeNames)-1])
	}
}
