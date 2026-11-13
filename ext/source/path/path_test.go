package path

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/Shougo/goda-cmp.vim/cmp"
)

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

func TestGatherDotFiles(t *testing.T) {
	dir := t.TempDir()
	for _, n := range []string{".env", "notes.txt"} {
		if err := os.WriteFile(filepath.Join(dir, n), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, tc := range []struct {
		name       string
		want       []string
		incomplete bool
	}{
		// Asked again, so that a dot typed next finds .env.
		{"", []string{"notes.txt"}, true},
		{".", []string{".env", "notes.txt"}, false},
		// A dot can no longer start the name.
		{"n", []string{"notes.txt"}, false},
	} {
		g, err := New().Gather(context.Background(), nil, &cmp.Context{Input: dir + "/" + tc.name})
		if err != nil {
			t.Fatal(err)
		}
		var got []string
		for _, it := range g.Items {
			got = append(got, it.Word)
		}
		slices.Sort(got)
		if !slices.Equal(got, tc.want) || g.Incomplete != tc.incomplete {
			t.Errorf("%q: got %v (incomplete %v), want %v (incomplete %v)",
				tc.name, got, g.Incomplete, tc.want, tc.incomplete)
		}
	}
}
