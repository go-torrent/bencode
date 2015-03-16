package bencode

import (
	"fmt"
	"sort"
	"strings"
)

type encodedPair struct {
	key   string
	value string
}

// implements sort.Interface
type byKey []*encodedPair

func (a byKey) Len() int           { return len(a) }
func (a byKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byKey) Less(i, j int) bool { return a[i].key < a[j].key }

func newEncodedPair(key, value interface{}) (*encodedPair, error) {
	encodedKey, err := Marshal(key)
	if err != nil {
		return nil, err
	}

	encodedValue, err := Marshal(value)
	if err != nil {
		return nil, err
	}

	return &encodedPair{string(encodedKey), string(encodedValue)}, nil
}

// Marshal returns the bencode encoding of v
func Marshal(v interface{}) ([]byte, error) {
	var encoded string

	switch t := v.(type) {
	case string:
		encoded = fmt.Sprintf("%d:%s", len(t), t)
	case int:
		encoded = fmt.Sprintf("i%de", t)
	case []interface{}:
		encodedElements := make([]string, len(t))

		for i, u := range t {
			encodedElement, err := Marshal(u)

			if err != nil {
				return nil, err
			}

			encodedElements[i] = string(encodedElement)
		}

		sort.Strings(encodedElements)
		encodedList := strings.Join(encodedElements, "")
		encoded = fmt.Sprintf("l%se", encodedList)
	case map[string]interface{}:
		encodedPairs := make([]*encodedPair, len(t))

		i := 0
		for k, v := range t {
			pair, err := newEncodedPair(k, v)
			if err != nil {
				return nil, err
			}

			encodedPairs[i] = pair
			i++
		}

		// it should make no difference to encode before sorting
		// because 3:*** is always < 4:****, for instance
		sort.Sort(byKey(encodedPairs))

		encodedDictionary := make([]string, len(encodedPairs)*2)
		for i, pair := range encodedPairs {
			encodedDictionary[2*i] = pair.key
			encodedDictionary[2*i+1] = pair.value
		}

		encoded = fmt.Sprintf("d%se", strings.Join(encodedDictionary, ""))
	default:
		// t isn't one of the types above
		return nil, fmt.Errorf("unsupported value %v", t)
	}

	return []byte(encoded), nil
}
