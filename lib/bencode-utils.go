package lib

import (
	"strconv"
	"strings"
	"unicode"
)

func determineNextElement(runes []rune, index int) (string, int) {
	if runes[index] == 'i' {
		return decodeNumber(runes, index)
	} else if runes[index] == 'd' {
		return decodeDict(runes, index)
	} else if runes[index] == 'l' {
		return decodeList(runes, index)
	} else if unicode.IsDigit(runes[index]) {
		return decodeStringLiteral(runes, index)
	} else {
		// error
	}
	return "", len(runes)
}

func decodeList(runes []rune, index int) (string, int) {
	index++
	builder := strings.Builder{}
	builder.WriteString("[")
	var str string
	for runes[index] != 'e' {
		str, index = determineNextElement(runes, index)
		builder.WriteString(str + ",")
	}
	temp := builder.String()
	temp = temp[:len(temp)-1] + "]"
	return temp, index + 1
}

func decodeNumber(runes []rune, index int) (string, int) {
	index++
	sign := 1
	result := 0
	if runes[index] == '-' {
		sign = -1
		index++
	}
	for runes[index] != 'e' {
		result = result*10 + int(runes[index]-'0')
		index++
	}
	return strconv.Itoa(result * sign), index + 1
}

func decodeStringLiteral(runes []rune, index int) (string, int) {
	strLen := 0
	for unicode.IsDigit(runes[index]) {
		strLen = strLen*10 + int(runes[index]-'0')
		index++
	}
	if runes[index] != ':' {
		// error
	}
	index++
	return string(runes[index : index+strLen]), index + strLen
}

func decodeDict(runes []rune, index int) (string, int) {
	builder := strings.Builder{}
	builder.WriteString("{")
	index++
	readingKey := true
	var str string
	for runes[index] != 'e' {
		if readingKey && unicode.IsDigit(runes[index]) {
			var str string
			str, index = decodeStringLiteral(runes, index)
			builder.WriteString(str + ":")
			readingKey = !readingKey
		} else if readingKey {
			// error
		} else if !readingKey {
			str, index = determineNextElement(runes, index)
			builder.WriteString(str + ",")
			readingKey = !readingKey
		} else {
			// error
		}
	}
	temp := builder.String()
	temp = temp[:len(temp)-1] + "}"
	return temp, index + 1
}

func BencodeToJSON(encodedStr string) string {
	runes := []rune(strings.ToLower(encodedStr))
	var decodedStr string
	if len(runes) > 0 {
		decodedStr, _ = determineNextElement(runes, 0)
		return decodedStr
	}
	return decodedStr
}
