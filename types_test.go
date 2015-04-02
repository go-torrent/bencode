package bencode

import (
	"testing"
)

func TestListAsStrings(t *testing.T) {
	expected := []string{"first", "second", "third"}
	l := List{
		"first", "second", "third",
	}

	actual, _ := l.AsStrings()

	if len(actual) != len(expected) {
		t.Fatalf("actual and expected have different sizes")
	}

	for k, v := range actual {
		if expected[k] != v {
			t.Errorf("expected\n%q got\n%q", expected[k], v)
		}
	}

}
