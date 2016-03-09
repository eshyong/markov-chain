package parse

import "testing"

func stringSlicesEqual(expected, actual []string) bool {
	for i, e := range expected {
		if e != actual[i] {
			return false
		}
	}
	return true
}

func TestGetNextUtf8Token(t *testing.T) {
	cases := map[string]string{
		"hi":            "hi",
		"hi123":         "hi123",
		"    hi":        "hi",
		"":              "",
		"hi\n":          "hi",
		"hi\t":          "hi",
		"hello_goodbye": "hello_goodbye",
	}
	for input, expected := range cases {
		actual, _ := getNextUtf8Token(input)
		if actual != expected {
			t.Errorf("Input: %q, Expected: %q, Actual: %q\n", input, expected, actual)
		}
	}
}

func TestParseInputString(t *testing.T) {
	cases := map[string][]string{
		"Parse. This!":                []string{"Parse", "This"},
		"Hello? Is anyone out there?": []string{"Hello", "Is", "anyone", "out", "there"},
		"Captain's log: day one.":     []string{"Captain's", "log", "day", "one"},
	}
	for input, expected := range cases {
		actual := ParseInputString(input)
		if !stringSlicesEqual(expected, actual) {
			t.Errorf("Input: %q, Expected: %q, Actual: %q\n", input, expected, actual)
		}
	}
}
