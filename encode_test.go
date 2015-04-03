package bencode

import (
	"testing"
)

func TestEncodeString(t *testing.T) {
	testCases := map[string]string{
		"some string": "11:some string",
		"":            "0:",
	}

	for value, expected := range testCases {
		encoded, err := Marshal(value)

		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}

		if string(encoded) != expected {
			t.Errorf("Marshal: expected %q, got %q", expected, encoded)
		}
	}
}

func TestEncodeInteger(t *testing.T) {
	testCases := map[string]int{
		"i-12e":        -12,
		"i0e":          0,
		"i123e":        123,
		"i5000000000e": 5e+9,
	}

	for expected, value := range testCases {
		encoded, err := Marshal(value)

		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}

		if expected != string(encoded) {
			t.Errorf("Marshal: expected %q, got %q", expected, encoded)
		}
	}
}

func TestEncodeList(t *testing.T) {
	testCases := map[string]List{
		"l4:spami42ee": List{"spam", 42},
		"le":           List{},
	}

	for expected, value := range testCases {
		encoded, err := Marshal(value)

		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}

		if string(encoded) != expected {
			t.Errorf("Marshal: expected %q, got %q", expected, encoded)
		}
	}
}

func TestEncodeDictionary(t *testing.T) {
	testCases := map[string]Dictionary{
		"d3:bar4:spam3:fooi42ee": Dictionary{
			"bar": "spam",
			"foo": 42,
		},
		"d8:announcei2e13:announce-listi1ee": Dictionary{
			"announce-list": 1,
			"announce":      2,
		},
		"de": Dictionary{},
	}

	for expected, value := range testCases {
		encoded, err := Marshal(value)

		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}

		if string(encoded) != expected {
			t.Errorf("Marshal: expected %q, got %q", expected, encoded)
		}
	}
}

func TestEncodeUnsupportedType(t *testing.T) {
	if _, err := Marshal(10.23); err == nil {
		t.Fatalf("Marshal: should error when type is unsupported")
	}
}

func TestEncodeListOfDictionaries(t *testing.T) {
	value := List{
		Dictionary{
			"length": 928670754,
			"path":   List{"Big_Buck_Bunny_1080p_surround_FrostWire.com.avi"},
		},
		Dictionary{
			"length": 5008,
			"path":   List{"PROMOTE_YOUR_CONTENT_ON_FROSTWIRE_01_06_09.txt"},
		},
		Dictionary{
			"length": 3456234,
			"path":   List{"Pressrelease_BickBuckBunny_premiere.pdf"},
		},
	}
	expected := "ld6:lengthi3456234e4:pathl39:Pressrelease_BickBuckBunny_premiere.pdfeed6:lengthi5008e4:pathl46:PROMOTE_YOUR_CONTENT_ON_FROSTWIRE_01_06_09.txteed6:lengthi928670754e4:pathl47:Big_Buck_Bunny_1080p_surround_FrostWire.com.avieee"
	encoded, err := Marshal(value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if string(encoded) != expected {
		t.Errorf("Marshal: expected %q, got %q", expected, encoded)
	}
}

func TestEncodeComplexValue(t *testing.T) {
	value := Dictionary{
		"number":  -42,
		"chinese": List{"你好", "中文"},
		"other": Dictionary{
			"foo":    List{1, 2, "yes", "no"},
			"foobar": Dictionary{},
			"bar": List{
				Dictionary{"1": "a", "2": "bb"},
				Dictionary{"1": "z", "2": "yy"},
			},
		},
	}

	expected := "d7:chinesel6:中文6:你好e6:numberi-42e5:otherd3:barld1:11:a1:22:bbed1:11:z1:22:yyee3:fool2:no3:yesi1ei2ee6:foobardeee"
	encoded, err := Marshal(value)

	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if string(encoded) != expected {
		t.Errorf("Marshal: expected %q, got %q", expected, encoded)
	}
}
