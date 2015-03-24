package bencode

// List is the definition of list according to bencode.
type List []interface{}

// Dictionary is the definition of list according to bencode.
type Dictionary map[string]interface{}

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
