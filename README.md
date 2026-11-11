# goda-cmp-ext-path.vim

[日本語](README.ja.md)

> [!WARNING]
> This is a beta (v0.0.x). goda-cmp.vim itself is under active development,
> and its bundled extensions are being moved into repositories of their own, so
> the API, the import paths and even the name of this repository may change.

A path source for [goda-cmp.vim](https://github.com/Shougo/goda-cmp.vim).

It offers the names in a directory as they are, and keeps the full path, the
modification time and whether it is a directory in the UserData of each
candidate (`path.Data`), so that a UI can show them.

- It completes only when the word before the cursor starts with `/`, `~/`,
  `./`, `../`, `$VAR/` or `${VAR}/`. A word such as `github.com/foo` is not a
  path.
- A relative path starts from the directory of the buffer.
- Dot files are offered only when a dot is typed.

## Differences from the bundled `file` source

The `file` source bundled with goda-cmp.vim offers only the rest of the name
after what was typed, and gives its candidates no full path. This source
offers the whole name, so that a fuzzy matcher can narrow the candidates down,
and a UI can read the full path and the modification time from `path.Data`.

## Installation

Install it with your plugin manager next to goda-cmp.vim. With
`g:goda#workspace` (goda-pm.vim sets it), `goda#build()` builds your
configuration with the installed plugin. Otherwise, add it to the Go module of
your configuration:

```sh
go get github.com/delphinus/goda-cmp-ext-path.vim
```

## Configuration

In the `main.go` of your configuration:

```go
import (
	"context"

	"github.com/Shougo/goda-cmp.vim/cmp"
	"github.com/Shougo/goda.vim/sdk"

	"github.com/delphinus/goda-cmp-ext-path.vim/ext/source/path"
)

func main() {
	b := cmp.New()
	// ...
	b.AddSource(path.New())
	b.PatchGlobal(cmp.Config{
		Sources: []string{"path", "around", "buffer"},
		// Offer only the candidates of path while a path is typed.
		DynamicSources: func(_ context.Context, _ sdk.Host, a cmp.DynamicSourcesArgs) ([]string, error) {
			if a.Context.Mode == "i" && path.InContext(a.Context.Input) {
				return []string{"path"}, nil
			}
			return nil, nil
		},
		SourceOptions: map[string]cmp.SourceOption{
			// Complete right after a / is typed.
			"path": {ForceCompletionPattern: path.ForceCompletionPattern},
		},
	})
	b.Main()
}
```

## License

MIT, see LICENSE.
