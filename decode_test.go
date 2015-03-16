package bencode

import (
	"testing"
)

func TestDecodeString(t *testing.T) {
	var value string
	expected := "some string"

	err := Unmarshal([]byte("11:some string"), &value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if value != expected {
		t.Errorf("Marshal: expected %v, got %v", expected, value)
	}
}

func TestDecodeInteger(t *testing.T) {
	var value int
	expected := -25

	err := Unmarshal([]byte("i-25e"), &value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if value != expected {
		t.Errorf("Marshal: expected %v, got %v", expected, value)
	}
}

func TestDecodeList(t *testing.T) {
	var value []interface{}
	expected := []interface{}{"spam", 42}

	err := Unmarshal([]byte("l4:spami42ee"), &value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if len(value) != len(expected) || value[0] != expected[0] || value[1] != value[1] {
		t.Errorf("Marshal: expected %v, got %v", expected, value)
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
		t.Fatalf("Marshal: %v", err)
	}

	if len(value) != len(expected) {
		t.Errorf("Marshal: expected %v, got %v", expected, value)
	}

	for k, v := range expected {
		if value[k] != v {
			t.Errorf("Marshal: expected %v, got %v", v, value[k])
		}
	}

}
