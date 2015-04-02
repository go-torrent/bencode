package bencode

import (
	"errors"
)

// List is the definition of list according to bencode.
type List []interface{}

// AsStrings returns the list as a list of strings
func (l List) AsStrings() ([]string, error) {
	strings := make([]string, len(l))

	for i, v := range l {
		switch s := v.(type) {
		case string:
			strings[i] = s
		default:
			return strings, errors.New("one element is not a string")
		}
	}

	return strings, nil
}

// Dictionary is the definition of list according to bencode.
type Dictionary map[string]interface{}

// Get returns the value of a key in the dictionary
func (d Dictionary) Get(key string) (interface{}, bool) {
	interfaceMap := map[string]interface{}(d)
	v, ok := interfaceMap[key]
	return v, ok
}

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
