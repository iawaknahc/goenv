package goenv

import (
	"reflect"
	"testing"
)

func TestParseEnvironment(t *testing.T) {
	input := []string{
		"a=",
		"b=1=2",
		"c=3",
	}
	actual := parseEnvironment(func() []string { return input })
	expected := environment{
		"a": "",
		"b": "1=2",
		"c": "3",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("%q != %q", actual, expected)
	}
}
