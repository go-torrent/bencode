package bencode

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

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
		value := List{}
		buffPos := 1

		for data[buffPos] != 'e' {
			decodedElement, pos, err := unmarshalFirstValue(data[buffPos:])
			if err != nil {
				return nil, buffPos + pos, err
			}

			//FIXME panic or error?
			if 0 > pos {
				panic(fmt.Sprintf("Failed to read bytes from %v", data[buffPos:]))
			}

			value = append(value, decodedElement)
			buffPos += pos

			if buffPos >= len(data) {
				panic("malformed value")
			}
		}

		// Includes the 'e' in the copied return
		return value, buffPos + 1, nil

	case dictionary:
		value := Dictionary{}
		buffPos := 1

		for data[buffPos] != 'e' {
			k, pos, err := unmarshalFirstValue(data[buffPos:])
			if err != nil {
				return nil, buffPos + pos, err
			}

			buffPos += pos

			v, pos, err := unmarshalFirstValue(data[buffPos:])
			if err != nil {
				return nil, buffPos + pos, err
			}

			// FIXME: improve type assertion
			strKey := k.(string)
			value[strKey] = v

			buffPos += pos
		}

		if buffPos >= len(data) {
			panic("malformed value")
		}

		// Includes the 'e' in the copied return
		return value, buffPos + 1, nil
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

		switch t := value.(type) {
		case Dictionary:
			rv.Set(reflect.ValueOf(t))
		default:
			return fmt.Errorf("decoded unexpected value %v of type %s", value, reflect.ValueOf(value).Type())
		}

	default:
		return fmt.Errorf("cant decode %v", data)
	}

	return nil
}
