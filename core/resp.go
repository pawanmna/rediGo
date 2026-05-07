// Package core contains core functitons of RESP
package core

import (
	"errors"
	"fmt"
	"strconv"
)

// // --- Length Parsing ---
// func readLength() {
// 	// TODO: parse digits until non-digit
// }
//

// DecodeArrayString is asserting that the inteface are string. (will cause problem if returned value from Decode is not string)
// the func is for Redis command frames only
func DecodeArrayString(data []byte) ([]string, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}
	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("expected RESP array")
	}
	tokens := make([]string, len(items))
	for i, item := range items {
		s, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("expected string at index %d, got %T", i, item)
		}
		tokens[i] = s
	}
	return tokens, nil
}

// --- Simple String (+) ---
func readSimpleString(data []byte) (any, int, error) {
	// TODO: skip '+', read until \r
	pos := 1
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}

	return string(data[1:pos]), pos + 2, nil
}

// --- Error (-) ---
func readError(data []byte) (any, int, error) {
	// TODO: reuse simple string logic
	return readSimpleString(data)
}

// --- Integer (:) ---
func readInt(data []byte) (any, int, error) {
	// TODO: parse integer until \r

	pos := 1
	for pos < len(data) && data[pos] != '\r' {
		pos++
	}

	v, err := strconv.ParseInt(string(data[1:pos]), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return v, pos + 2, nil
}

// --- Bulk String ($) ---
func readBulkString(data []byte) (any, int, error) {
	// TODO:
	x, pos, err := readInt(data)
	if err != nil {
		return nil, 0, err
	}
	length := int(x.(int64))

	return string(data[pos:(pos + length)]), length + pos + 2, nil
	// 1. read length
	// 2. advance position
	// 3. read `length` bytes
}

// --- Array (*) ---
func readArray(data []byte) (any, int, error) {
	// TODO:
	// 1. read element count
	// 2. loop and decode each element
	x, pos, err := readInt(data)
	if err != nil {
		return nil, 0, err
	}
	size := int(x.(int64))
	elems := make([]any, size)
	for i := range elems {
		elem, del, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += del
	}

	return elems, pos, nil
}

func DecodeOne(data []byte) (any, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}

	switch data[0] {
	case '+':
		return readSimpleString(data)
	case ':':
		return readInt(data)
	case '$':
		return readBulkString(data)
	case '*':
		return readArray(data)
	case '-':
		return readError(data)
	}

	return nil, 0, fmt.Errorf("unknown RESP type: %q", data[0])
	// TODO:
	// 1. check empty input
	// 2. switch on first byte
	// 3. delegate to correct reader
}

func Decode(data []byte) (any, error) {
	// TODO:
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	val, _, err := DecodeOne(data)
	return val, err
	// 1. validate input
	// 2. call DecodeOne
}

func Encode(arg any, isSimple bool) ([]byte, error) {
	switch v := arg.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v)), nil
		}
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)), nil
	}
	return []byte{}, errors.New("empty")
}
