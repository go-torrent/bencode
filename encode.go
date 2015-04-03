package bencode

import (
	"fmt"
	"sort"
	"strings"
)

func encodePair(key, value interface{}) ([]byte, []byte, error) {
	encodedKey, err := Marshal(key)
	if err != nil {
		return nil, nil, err
	}

	encodedValue, err := Marshal(value)
	if err != nil {
		return nil, nil, err
	}

	return encodedKey, encodedValue, nil
}

// Marshal returns the bencode encoding of v
func Marshal(v interface{}) ([]byte, error) {
	var encoded string

	switch t := v.(type) {
	case string:
		encoded = fmt.Sprintf("%d:%s", len(t), t)

	case int:
		encoded = fmt.Sprintf("i%de", t)

	case List:
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

	case Dictionary:
		sortedKeys := make([]string, 0, len(t))
		for k := range t {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Sort(sort.StringSlice(sortedKeys))

		encodedDictionary := make([]string, 0, len(t)*2)
		for _, k := range sortedKeys {
			encK, encV, err := encodePair(k, t[k])
			if err != nil {
				return nil, err
			}
			encodedDictionary = append(encodedDictionary, string(encK), string(encV))
		}

		//FIXME this can probably be optimized
		encoded = fmt.Sprintf("d%se", strings.Join(encodedDictionary, ""))

	default:
		return nil, fmt.Errorf("unsupported value %v", t)
	}

	return []byte(encoded), nil
}
