package bencode

import (
	"testing"
)

func TestUnmarshalFirstValue(t *testing.T) {
	testCases := map[string]int{
		"11:some string":                     14,
		"i-25e":                              5,
		"l4:spami42ee":                       12,
		"d3:bar4:spam3:fooi42ee":             22,
		"ld1:a2:aa1:b2:bbed1:c2:cc1:d2:ddee": 34,
		"d1:a2:aae_________":                 9,
		"l1:24:5678e_-_-_-_":                 11,
	}

	for encoded, size := range testCases {
		decoded, pos, _ := unmarshalFirstValue([]byte(encoded))

		t.Logf("unmarshalFirstValue: %q => %v\n", encoded, decoded)

		if pos != size {
			t.Fatalf("unmarshalFirstValue: expected\n %v, got\n %v", size, pos)
		}
	}
}

func TestDecodeString(t *testing.T) {
	var value string
	expected := "some string"

	err := Unmarshal([]byte("11:some string"), &value)

	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if value != expected {
		t.Errorf("Unmarshal: expected %v, got %v", expected, value)
	}
}

func TestDecodeInteger(t *testing.T) {
	var value int
	expected := -25

	err := Unmarshal([]byte("i-25e"), &value)

	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if value != expected {
		t.Errorf("Unmarshal: expected %v, got %v", expected, value)
	}
}

func TestDecodeList(t *testing.T) {
	var value []interface{}
	expected := []interface{}{"spam", 42}

	err := Unmarshal([]byte("l4:spami42ee"), &value)

	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(value) != len(expected) || value[0] != expected[0] || value[1] != value[1] {
		t.Errorf("Unmarshal: expected %v, got %v", expected, value)
	}
}

func TestDecodeDictionary(t *testing.T) {
	var value map[string]interface{}
	expected := map[string]interface{}{
		"bar": "spam",
		"foo": 42,
	}

	err := Unmarshal([]byte("d3:bar4:spam3:fooi42ee"), &value)

	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(value) != len(expected) {
		t.Errorf("Unmarshal: expected %v, got %v", expected, value)
	}

	for k, v := range expected {
		if value[k] != v {
			t.Errorf("Unmarshal: expected %v, got %v", v, value[k])
		}
	}

}
