package assert

import (
	"encoding/json"
	"math/big"
	"strings"
)

func decodeJSONNumbers(text string) (interface{}, error) {
	decoder := json.NewDecoder(strings.NewReader(text))
	decoder.UseNumber()
	var value interface{}
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	return value, nil
}

func normalizeJSONNumbers(value interface{}) interface{} {
	switch value := value.(type) {
	case json.Number:
		return normalizeJSONNumber(value)
	case []interface{}:
		for i, element := range value {
			value[i] = normalizeJSONNumbers(element)
		}
	case map[string]interface{}:
		for key, element := range value {
			value[key] = normalizeJSONNumbers(element)
		}
	}
	return value
}

func normalizeJSONNumber(value json.Number) json.Number {
	number := string(value)
	negative := strings.HasPrefix(number, "-")
	if negative {
		number = number[1:]
	}
	exponentText := "0"
	if i := strings.IndexAny(number, "eE"); i >= 0 {
		exponentText = number[i+1:]
		number = number[:i]
	}
	var fraction int
	if i := strings.IndexByte(number, '.'); i >= 0 {
		fraction = len(number) - i - 1
		number = number[:i] + number[i+1:]
	}
	number = strings.TrimLeft(number, "0")
	if number == "" {
		return "0"
	}
	digits := strings.TrimRight(number, "0")
	// Keep the exponent separate so tiny numbers need no enormous power of ten.
	var exponent big.Int
	exponent.SetString(exponentText, 10)
	exponent.Add(&exponent, big.NewInt(int64(len(number)-len(digits)-fraction)))
	if negative {
		digits = "-" + digits
	}
	return json.Number(digits + "e" + exponent.String())
}
