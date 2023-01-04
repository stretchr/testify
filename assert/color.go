package assert

import "fmt"

type color string

const (
	resetCode     color = "\033[0m"
	redCode       color = "\033[31m"
	brightRedCode color = "\033[31;1m"
	greenCode     color = "\033[32m"
	yellowCode    color = "\033[33m"
	blueCode      color = "\033[34m"

	setFormat = "%s%s%s"
)

// formatStr formats string by argument placing
func formatStr(f string, a ...interface{}) string {
	return fmt.Sprintf(f, a...)
}

// red prints text in red
func red(format string, args ...interface{}) string {
	return fmt.Sprintf(setFormat, redCode, formatStr(format, args...), resetCode)
}

// brightRed prints text in red, brighter and thicker
func brightRed(format string, args ...interface{}) string {
	return fmt.Sprintf(setFormat, brightRedCode, formatStr(format, args...), resetCode)
}

// green prints text in green
func green(format string, args ...interface{}) string {
	return fmt.Sprintf(setFormat, greenCode, formatStr(format, args...), resetCode)
}

// yellow prints text in yellow
func yellow(format string, args ...interface{}) string {
	return fmt.Sprintf(setFormat, yellowCode, formatStr(format, args...), resetCode)
}

// blue prints text in blue
func blue(format string, args ...interface{}) string {
	return fmt.Sprintf(setFormat, blueCode, formatStr(format, args...), resetCode)
}
