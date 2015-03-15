package bencode

import (
	"testing"
)

func TestEncodeString(t *testing.T) {
	value := "some string"
	expected := "11:some string"
	encoded, err := Marshal(value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if string(encoded) != expected {
		t.Errorf("Marshal: expected %q, got %q", expected, encoded)
	}
}

func TestEncodeInteger(t *testing.T) {
	value := -12
	expected := "i-12e"
	encoded, err := Marshal(value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if string(encoded) != expected {
		t.Errorf("Marshal: expected %q, got %q", expected, encoded)
	}
}

func TestEncodeList(t *testing.T) {
	value := []interface{}{"spam", 42}
	expected := "l4:spami42ee"
	encoded, err := Marshal(value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if string(encoded) != expected {
		t.Errorf("Marshal: expected %q, got %q", expected, encoded)
	}
}

func TestEncodeDictionary(t *testing.T) {
	value := map[string]interface{}{
		"bar": "spam",
		"foo": 42,
	}
	expected := "d3:bar4:spam3:fooi42ee"
	encoded, err := Marshal(value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if string(encoded) != expected {
		t.Errorf("Marshal: expected %q, got %q", expected, encoded)
	}
}

func TestEncodeUnsupportedType(t *testing.T) {
	if _, err := Marshal(10.23); err == nil {
		t.Fatalf("Marshal: should error when type is unsupported")
	}
}
