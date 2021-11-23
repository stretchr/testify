package assert

import "regexp"

const (
	expected_colon = "expected:"
	actual_colon   = "actual  :"
)

func colorize(plainText string) string {
	re := regexp.MustCompile(`(.*)\n` + expected_colon + `(.*)\n` + actual_colon + `(.*)`)

	return string(
		re.ReplaceAll(
			[]byte(plainText),
			[]byte("\033[33m${1}\033[0m\n"+
				expected_colon+"\033[32m${2}\033[0m\n"+
				actual_colon+"\033[31m${3}\033[0m"),
		),
	)
}
