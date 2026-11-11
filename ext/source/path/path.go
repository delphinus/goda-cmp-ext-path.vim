// Package path is a source "path" for goda-cmp.vim that completes file
// paths. It behaves like the path source of blink.cmp.
//
// The file source of goda-cmp-ext gives its candidates no full path, so
// a UI cannot show their modification times, and it offers only the rest of
// the name after what was typed, so a fuzzy matcher cannot narrow them down.
// This source offers the names in the directory as they are, and puts the
// full path, the modification time and whether it is a directory in the
// UserData (Data). A UI that reads it can show, for example, the
// modification time of each candidate.
//
// It completes only when the word before the cursor starts with / ~/ ./ ../
// $VAR/ or ${VAR}/ (a quote or a ( right before it is allowed). A word such
// as github.com/foo is not a path. A relative path starts from the
// directory of the buffer (the default of blink.cmp). Dot files are offered
// only when a dot is typed.
package path

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Shougo/goda-cmp.vim/cmp"
	"github.com/Shougo/goda.vim/sdk"
)

// Data is the UserData of a candidate.
type Data struct {
	Full  string
	MTime time.Time
	Dir   bool
}

// ForceCompletionPattern is the source option (ForceCompletionPattern) that
// starts the completion right after a / is typed. Gather decides whether the
// word is really a path.
const ForceCompletionPattern = `(?:^|[\s"'\x60(])(?:~|\.\.?|\$\{?\w+\}?)?/\S*`

var (
	token  = regexp.MustCompile(`\S*$`)
	marker = regexp.MustCompile(`^(?:/|~/|\.\.?/|\$\w+/|\$\{\w+\}/)`)
	envVar = regexp.MustCompile(`^\$\{?(\w+)\}?`)
)

type source struct{}

// New returns the source "path".
func New() cmp.Source { return source{} }

func (source) Name() string { return "path" }

// typed returns the path before the cursor and the byte offset where it starts.
func typed(input string) (string, int, bool) {
	loc := token.FindStringIndex(input)
	start := loc[0]
	t := input[start:]
	if t != "" && strings.ContainsRune(`"'`+"`(", rune(t[0])) {
		t, start = t[1:], start+1
	}
	if !marker.MatchString(t) {
		return "", 0, false
	}
	return t, start, true
}

// InContext reports whether the word before the cursor is a path, by the same
// rule as Gather. Use it, for example in DynamicSources, to offer only the
// candidates of "path" while a path is being typed.
func InContext(input string) bool {
	_, _, ok := typed(input)
	return ok
}

// CompletePosition is the position after the last / (where the name to
// complete starts).
func (source) CompletePosition(_ context.Context, _ sdk.Host, c *cmp.Context) int {
	t, start, ok := typed(c.Input)
	if !ok {
		return -1
	}
	return start + strings.LastIndexByte(t, '/') + 1
}

func (s source) Gather(ctx context.Context, h sdk.Host, c *cmp.Context) (cmp.Gathered, error) {
	t, start, ok := typed(c.Input)
	if !ok {
		return cmp.Gathered{Start: len(c.Input)}, nil
	}
	slash := strings.LastIndexByte(t, '/')
	dir, name := t[:slash+1], t[slash+1:]
	g := cmp.Gathered{Start: start + slash + 1}

	switch {
	case strings.HasPrefix(dir, "~/"):
		home, err := os.UserHomeDir()
		if err != nil {
			return g, err
		}
		dir = filepath.Join(home, dir[2:])
	case strings.HasPrefix(dir, "$"):
		m := envVar.FindStringSubmatch(dir)
		dir = os.Getenv(m[1]) + dir[len(m[0]):]
	case strings.HasPrefix(dir, "."):
		base, err := h.Call(ctx, "expand", "#"+strconv.Itoa(c.Bufnr)+":p:h")
		if err != nil {
			return g, err
		}
		b, _ := base.(string)
		if b == "" {
			b, _ = os.Getwd()
		}
		dir = filepath.Join(b, dir)
	}

	ents, err := os.ReadDir(dir)
	if err != nil {
		return g, nil
	}
	hidden := strings.HasPrefix(name, ".")
	for _, e := range ents {
		n := e.Name()
		if strings.HasPrefix(n, ".") && !hidden {
			continue
		}
		full := filepath.Join(dir, n)
		st, err := os.Stat(full) // follow symbolic links
		if err != nil {
			continue
		}
		it := cmp.Item{Word: n, Kind: "File"}
		d := Data{Full: full, MTime: st.ModTime(), Dir: st.IsDir()}
		if d.Dir {
			it.Word += "/"
			it.Kind = "Folder"
		}
		it.UserData = d
		g.Items = append(g.Items, it)
		if ctx.Err() != nil {
			return g, ctx.Err()
		}
	}
	return g, nil
}

var (
	_ cmp.Source             = source{}
	_ cmp.CompletePositioner = source{}
)
