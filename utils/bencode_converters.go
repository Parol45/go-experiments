package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func byteIsDigit(symbol byte) bool {
	return 47 < symbol && symbol < 58
}

func tryDecodeIpOrHash(bytes []byte, key string) string {
	if key == "\"ip\"" {
		ip1, err1 := strconv.ParseInt(fmt.Sprintf("%x", bytes[0]), 16, 64)
		ip2, err2 := strconv.ParseInt(fmt.Sprintf("%x", bytes[1]), 16, 64)
		ip3, err3 := strconv.ParseInt(fmt.Sprintf("%x", bytes[2]), 16, 64)
		ip4, err4 := strconv.ParseInt(fmt.Sprintf("%x", bytes[3]), 16, 64)
		port, err5 := strconv.ParseInt(fmt.Sprintf("%x", bytes[4])+fmt.Sprintf("%x", bytes[5]), 16, 64)
		if err1 == nil && err2 == nil && err3 == nil && err4 == nil && err5 == nil {
			return fmt.Sprintf("%d.%d.%d.%d:%d", ip1, ip2, ip3, ip4, port)
		}
	}
	var result strings.Builder
	for _, b := range bytes {
		str := fmt.Sprintf("%02x", b)
		result.WriteString(str)
	}
	return result.String()
}

func decodeNextElement(bytes []byte, index int, key string) (string, int, error) {
	if bytes[index] == 'i' {
		return decodeNumber(bytes, index)
	} else if bytes[index] == 'd' {
		return decodeDict(bytes, index)
	} else if bytes[index] == 'l' {
		return decodeList(bytes, index)
	} else if byteIsDigit(bytes[index]) {
		return decodeStringLiteral(bytes, index, key)
	} else {
		return "", len(bytes), errors.New(fmt.Sprintf("No known entity start for index: %d, symbol: '%s'", index, string(bytes[index])))
	}
}

func decodeList(bytes []byte, index int) (string, int, error) {
	startIndex := index
	index++
	builder := strings.Builder{}
	builder.WriteString("[")
	var tempStr string
	var err error
	for index < len(bytes) && bytes[index] != 'e' {
		tempStr, index, err = decodeNextElement(bytes, index, "")
		if err != nil {
			return "", 0, err
		}
		builder.WriteString(tempStr + ",")
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for list starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	temp := builder.String()
	temp = temp[:len(temp)-1] + "]"
	return temp, index + 1, err
}

func decodeNumber(bytes []byte, index int) (string, int, error) {
	startIndex := index
	index++
	result := strings.Builder{}
	for index < len(bytes) && bytes[index] != 'e' {
		result.WriteByte(bytes[index])
		index++
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for number starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	return result.String(), index + 1, nil
}

func decodeStringLiteral(bytes []byte, index int, key string) (string, int, error) {
	strLen := 0
	for byteIsDigit(bytes[index]) {
		strLen = strLen*10 + int(bytes[index]-'0')
		index++
	}
	if bytes[index] != ':' {
		return "", 0, errors.New(fmt.Sprintf("Missing string semicolon. Index: %d, Symbol: '%s'", index, string(bytes[index])))
	} else if index+strLen > len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("Wrong string length. Index: %d, Symbol: '%s'", index, string(bytes[index])))
	}
	index++
	var resultingStr string
	if !InStringArray(key, []string{"\"y\"", "\"q\"", ""}) {
		resultingStr = "\"" + tryDecodeIpOrHash(bytes[index:index+strLen], key) + "\""
	} else {
		resultingStr = "\"" + string(bytes[index:index+strLen]) + "\""
	}
	return resultingStr, index + strLen, nil
}

func decodeDict(bytes []byte, index int) (string, int, error) {
	startIndex := index
	builder := strings.Builder{}
	builder.WriteString("{")
	index++
	readingKey := true
	var tempStr, key string
	var err error
	for index < len(bytes) && bytes[index] != 'e' {
		if readingKey && byteIsDigit(bytes[index]) {
			key, index, err = decodeStringLiteral(bytes, index, "")
			if err != nil {
				return "", 0, err
			}
			builder.WriteString(key + ":")
			readingKey = !readingKey
		} else if readingKey {
			return "", 0, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(bytes[index])))
		} else if !readingKey {
			tempStr, index, err = decodeNextElement(bytes, index, key)
			if err != nil {
				return "", 0, err
			}
			builder.WriteString(tempStr + ",")
			readingKey = !readingKey
		}
	}
	if index >= len(bytes) {
		return "", 0, errors.New(fmt.Sprintf("No closing 'e' symbol for dictionary starting at index: %d, symbol: '%s'", startIndex, string(bytes[startIndex])))
	}
	temp := builder.String()
	temp = temp[:len(temp)-1] + "}"
	return temp, index + 1, err
}

func encodeNextElement(json string, index int, key string) ([]byte, int, error) {
	if byteIsDigit(json[index]) {
		return encodeNumber(json, index)
	} else if json[index] == '{' {
		return encodeDict(json, index)
	} else if json[index] == '[' {
		return encodeList(json, index)
	} else if json[index] == '"' {
		return encodeStringLiteral(json, index, key)
	} else {
		return nil, 0, errors.New(fmt.Sprintf("No known entity start for index: %d, symbol: '%s'", index, string(json[index])))
	}
}

func encodeIpOrHash(literal string, key string) ([]byte, error) {
	var result []byte
	if key == "ip" {
		parts := strings.Split(strings.Replace(literal, ":", ".", -1), ".")
		i := 0
		for i < 4 {
			ipPart := parts[i]
			ipInt, err := strconv.Atoi(ipPart) // ahoy
			if err != nil || ipInt > 255 || ipInt < 0 {
				return nil, errors.New(fmt.Sprintf("Problems with parsing ip: '%s', part: '%s'", literal, ipPart))
			}
			result = append(result, byte(ipInt))
			i++
		}
		ipPart := parts[4]
		portInt, err := strconv.Atoi(ipPart) // ahoy
		if err != nil || portInt > 65535 {
			return nil, errors.New(fmt.Sprintf("Problems with parsing ip: '%s', part: '%s'", literal, ipPart))
		}
		portStr := fmt.Sprintf("%x", portInt)
		first, _ := hex.DecodeString(portStr[:2])
		second, _ := hex.DecodeString(portStr[2:])
		result = append(result, first...)
		result = append(result, second...)
		return result, nil
	}
	var i int
	if len(literal)%2 != 0 {
		return nil, errors.New(fmt.Sprintf("Even number of ascii symbols while converting to []bytes string: '%s'", literal))
	}
	for 2*i+1 < len(literal) {
		byteAt0, err := hex.DecodeString(literal[2*i : 2*i+2])
		if err != nil || len(byteAt0) > 1 || byteAt0[0] > 255 {
			return nil, errors.New(fmt.Sprintf("Problem with encoding hex string: '%s'", literal))
		}
		result = append(result, byteAt0[0])
		i++
	}
	return result, nil
}

func encodeDict(json string, index int) ([]byte, int, error) {
	startIndex := index
	index++
	readingKey := true
	var tempBytes, resultingBytes []byte
	var err error
	var key string
	resultingBytes = append(resultingBytes, 'd')
	for index < len(json) && json[index] != '}' {
		if readingKey && json[index] == ':' {
			readingKey = !readingKey
		} else if readingKey && json[index] == '"' {
			tempBytes, index, err = encodeStringLiteral(json, index, "")
			if err != nil {
				return nil, 0, err
			}
			newKey := string(tempBytes)
			newKey = newKey[strings.IndexByte(newKey, ':')+1:]
			if key != "" && strings.Compare(key, newKey) == 1 {
				return nil, 0, errors.New(fmt.Sprintf("No alphabetical order in dictionary key at index: %d", index))
			}
			key = newKey
			resultingBytes = append(resultingBytes, tempBytes...)
		} else if readingKey && json[index] != ',' {
			return nil, 0, errors.New(fmt.Sprintf("Dictionary key can be nothing but string. Index: %d, Symbol: '%s'", index, string(json[index])))
		} else if !readingKey {
			tempBytes, index, err = encodeNextElement(json, index, key)
			if err != nil {
				return nil, 0, err
			}
			resultingBytes = append(resultingBytes, tempBytes...)
			readingKey = !readingKey
		}
		index++
	}
	if index >= len(json) {
		return nil, 0, errors.New(fmt.Sprintf("No closing '}' symbol for dictionary starting at index: %d", startIndex))
	}
	resultingBytes = append(resultingBytes, 'e')
	return resultingBytes, index, err
}

func encodeList(json string, index int) ([]byte, int, error) {
	startIndex := index
	index++
	var tempBytes, resultingBytes []byte
	var err error
	resultingBytes = append(resultingBytes, 'l')
	for index < len(json) && json[index] != ']' {
		if json[index] != ',' {
			tempBytes, index, err = encodeNextElement(json, index, "")
			resultingBytes = append(resultingBytes, tempBytes...)
			if err != nil {
				return nil, 0, err
			}
		}
		index++
	}
	if index >= len(json) {
		return nil, 0, errors.New(fmt.Sprintf("No closing ']' symbol for list starting at index: %d", startIndex))
	}
	resultingBytes = append(resultingBytes, 'e')
	return tempBytes, index + 1, err
}

func encodeStringLiteral(json string, index int, key string) ([]byte, int, error) {
	index++
	literalEnd := strings.IndexByte(json[index:], '"')
	literal := json[index : index+literalEnd]
	var resultingBytes, temp []byte
	var err error
	var length int
	if !InStringArray(key, []string{"y", "q", ""}) {
		temp, err = encodeIpOrHash(literal, key)
		length = len(temp)
		resultingBytes = append(resultingBytes, fmt.Sprintf("%d:", length)...)
		resultingBytes = append(resultingBytes, temp...)
		if err != nil {
			return nil, 0, err
		}
	} else {
		length = len(literal)
		resultingBytes = append(resultingBytes, fmt.Sprintf("%d:", length)...)
		resultingBytes = append(resultingBytes, []byte(literal)...)
	}
	return resultingBytes, index + literalEnd, nil
}

func encodeNumber(json string, index int) ([]byte, int, error) {
	result := strings.Builder{}
	result.WriteByte('i')
	for index < len(json) && (byteIsDigit(json[index]) || json[index] == '-') {
		result.WriteByte(json[index])
		index++
	}
	if index < len(json) && InByteArray(json[index], []byte{'.', ','}) {
		return nil, 0, errors.New(fmt.Sprintf("Numbers must be integer. Index: %d", index))
	}
	result.WriteByte('e')
	return []byte(result.String()), index + 1, nil
}

// a - arguments (string)
// q - method name (string)
// y - message type: r - response, q - query, e - error (string)
// v - version (hex)
// t - transaction id (hex)
// id - id of node (hex)
// ip - ... (hex)
func BencodeToJSON(encodedStr []byte) (string, error) {
	var decodedStr string
	var err error
	if len(encodedStr) > 0 {
		decodedStr, _, err = decodeNextElement(encodedStr, 0, "")
		return decodedStr, err
	}
	return decodedStr, err
}

func JSONToBencode(json string) ([]byte, error) {
	var formattedStr string
	formattedStr = strings.Replace(json, " ", "", -1)
	formattedStr = strings.Replace(json, "\n", "", -1)
	formattedStr = strings.Replace(json, "\r", "", -1)
	formattedStr = strings.Replace(json, "\t", "", -1)
	var encodedBytes []byte
	var err error
	if len(formattedStr) > 0 {
		encodedBytes, _, err = encodeNextElement(formattedStr, 0, "")
		return encodedBytes, err
	}
	return encodedBytes, err
}
