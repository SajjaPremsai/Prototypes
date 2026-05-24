package core

import (
	"errors"
	"fmt"
)

func readSimpleString(data []byte) (string, int, error) {
	pos := 1
	for ; pos < len(data) && data[pos] != '\r'; pos++ {
	}

	if pos+1 >= len(data) {
		return "", 0, errors.New("malformed simple string")
	}

	return string(data[1:pos]), pos + 2, nil
}

func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

func readInteger(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	for ; pos < len(data) && data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}

	if pos+1 >= len(data) {
		return 0, 0, errors.New("malformed integer")
	}

	return value, pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	length, delta := readLength(data, pos)
	if delta < 0 {
		return "", 0, errors.New("malformed bulk string length")
	}

	pos += delta
	if pos+length+2 > len(data) {
		return "", 0, errors.New("malformed bulk string data")
	}

	return string(data[pos : pos+length]), pos + length + 2, nil
}

func readLength(data []byte, pos int) (int, int) {
	length := 0
	start := pos
	for pos < len(data) && data[pos] != '\r' {
		b := data[pos]
		if b < '0' || b > '9' {
			return 0, -1
		}
		length = length*10 + int(b-'0')
		pos++
	}

	if pos+1 >= len(data) || data[pos] != '\r' || data[pos+1] != '\n' {
		return 0, -1
	}

	return length, (pos - start) + 2
}

func readArray(data []byte) (interface{}, int, error) {
	pos := 1

	count, delta := readLength(data, pos)
	if delta < 0 {
		return nil, 0, errors.New("malformed array length")
	}
	pos += delta

	elems := make([]interface{}, count)
	for i := range elems {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("Empty data")
	}
	value, _, err := DecodeOne(data)
	return value, err
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("Empty data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInteger(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	}

	return nil, 0, nil
}

func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	ts := value.([]interface{})
	tokens := make([]string, len(ts))
	for i := range tokens {
		tokens[i] = ts[i].(string)
	}
	return tokens, nil
}

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimpleString {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
	}
	return []byte{}
}