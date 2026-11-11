package path

import "testing"

func TestInContext(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  bool
	}{
		{"./", true},
		{"see ~/doc", true},
		{`"../x`, true},
		// Only a quote or a ( at the very start of the word is skipped.
		{`open("../x`, false},
		{"$HOME/", true},
		{"github.com/foo", false},
		{"a/b", false},
		{"word", false},
	} {
		if got := InContext(tc.input); got != tc.want {
			t.Errorf("%q: got %v, want %v", tc.input, got, tc.want)
		}
	}
}
