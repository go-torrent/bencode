package bencode

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Type is the representation of a bencode type.
type Type uint

// The zero type is not a valid type.
const (
	invalid Type = iota
	integer
	str
	list
	dictionary
)

var typeNames = []string{
	invalid:    "invalid",
	integer:    "integer",
	str:        "string",
	list:       "list",
	dictionary: "dictionary",
}

func (t Type) String() string {
	return typeNames[t]
}

// TypeOf returns the type of the bencode-encoded data
func TypeOf(data []byte) Type {
	t := data[0]

	// FIXME: add additional validation

	switch {
	default:
		return invalid
	case t == 'i':
		return integer
	case t >= '0' && t <= '9':
		return str
	case t == 'l':
		return list
	case t == 'd':
		return dictionary
	}
}

func unmarshalFirstValue(data []byte) (value interface{}, copied int, err error) {
	switch TypeOf(data) {
	default:
		return nil, -1, fmt.Errorf("invalid data %q", data)

	case integer:
		endPos := strings.IndexRune(string(data), 'e')
		if endPos == -1 {
			return nil, 0, fmt.Errorf("invalid data %s", data)
		}

		//FIXME: leading zeros are not allowed
		//FIXME: minus 0 is not allowed

		value, err := strconv.Atoi(string(data)[1:endPos])
		if err != nil {
			return nil, endPos, err
		}

		// Includes the 'e' in the copied return
		return value, endPos + 1, nil

	case str:
		sepPos := strings.IndexRune(string(data), ':')
		if sepPos == -1 {
			return nil, sepPos, fmt.Errorf("invalid data %s", data)
		}

		size, err := strconv.Atoi(string(data)[0:sepPos])
		if err != nil {
			return nil, sepPos, err
		}

		//FIXME: check data len
		sepPos++
		endPos := sepPos + size
		value := string(data[sepPos:endPos])

		return value, endPos, nil

	case list:
		value := []interface{}{}
		for strPos := 1; strPos < len(data) && data[strPos] != 'e'; {
			decodedElement, pos, err := unmarshalFirstValue(data[strPos:])
			if err != nil {
				return nil, strPos + pos, err
			}

			//FIXME panic or error?
			if 0 > pos {
				panic(fmt.Sprintf("Failed to read bytes from %v", data[strPos:]))
			}

			value = append(value, decodedElement)
			strPos += pos
		}

		//FIXME check if read all the data

		return value, len(data), nil

	case dictionary:
		value := map[string]interface{}{}

		for strPos := 1; strPos < len(data) && data[strPos] != 'e'; {
			k, pos, err := unmarshalFirstValue(data[strPos:])
			if err != nil {
				return nil, strPos + pos, err
			}

			strPos += pos

			v, pos, err := unmarshalFirstValue(data[strPos:])
			if err != nil {
				return nil, strPos + pos, err
			}

			// FIXME: improve type assertion
			strKey := k.(string)
			value[strKey] = v

			strPos += pos
		}

		return value, len(data), nil
	}
}

// Unmarshal parses the bencode-encoded data and stores the result in the value
// pointed to by v
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("can't unmarshal to a %s variable", reflect.TypeOf(v))
	}

	rv = rv.Elem()

	switch TypeOf(data) {
	case integer:
		if rv.Kind() != reflect.Int {
			return fmt.Errorf("can't unmarshal int to a %s variable", rv.Kind())
		}

		//FIXME: leading zeros are not allowed
		//FIXME: minus 0 is not allowed

		value, _, err := unmarshalFirstValue(data)
		if err != nil {
			return err
		}

		//FIXME: check if value is Integer

		rv.Set(reflect.ValueOf(value))

	case str:
		if rv.Kind() != reflect.String {
			return fmt.Errorf("can't unmarshal string to a %s variable", rv.Kind())
		}

		value, _, err := unmarshalFirstValue(data)
		if err != nil {
			return err
		}

		//FIXME: check if value is String

		rv.Set(reflect.ValueOf(value))

	case list:
		switch rv.Kind() {
		default:
			return fmt.Errorf("can't unmarshal list to a %s variable", rv.Kind())
		case reflect.Array, reflect.Slice:
		}

		value, _, err := unmarshalFirstValue(data)
		if err != nil {
			return err
		}

		//FIXME: check if value is Array or Slice

		rv.Set(reflect.ValueOf(value))

	//FIXME: could it be a Struct rather than a map[string]T ?
	case dictionary:
		if destKind := rv.Kind(); destKind != reflect.Map {
			return fmt.Errorf("can't unmarshal dictionary to a %s variable", destKind)
		}

		t := rv.Type()
		if keyKind := t.Key().Kind(); keyKind != reflect.String {
			return fmt.Errorf("map key has wrong type: expected string, got %s", keyKind)
		}

		value, _, err := unmarshalFirstValue(data)
		if err != nil {
			return err
		}

		//FIXME: check if value is map

		rv.Set(reflect.ValueOf(value))

	default:
		//nothing
	}

	return nil
}
