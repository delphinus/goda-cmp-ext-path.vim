# goda-cmp-ext-path.vim

[日本語](README.ja.md)

> [!NOTE]
> goda-cmp.vim itself is under active development, so the API of this plugin
> may change to follow it.

https://github.com/user-attachments/assets/9c510911-74ec-4d9f-b9dc-0944b1c3951b

*The whole names with their modification times, narrowed down by `matcherfuzzy` of [goda-cmp-ext-fuzzy.vim](https://github.com/delphinus/goda-cmp-ext-fuzzy.vim); the dot files come when a dot is typed. The menu shows `path.Data` through a converter of the demo (recorded with `demo/record.sh`).*

A path source for [goda-cmp.vim](https://github.com/Shougo/goda-cmp.vim).

It offers the names in a directory as they are, and keeps the full path, the
modification time and whether it is a directory in the UserData of each
candidate (`path.Data`), so that a UI can show them.

- It completes only when the word before the cursor starts with `/`, `~/`,
  `./`, `../`, `$VAR/` or `${VAR}/`. A word such as `github.com/foo` is not a
  path.
- A relative path starts from the directory of the buffer.
- Dot files are offered only when a dot is typed.

## Differences from the `file` source of goda-cmp-ext

The `file` source of [goda-cmp-ext](https://github.com/Shougo/goda-cmp-ext) offers only the rest of the name
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

## Demo

`demo/record.sh` records the demo with a minimal configuration
(`demo/init.lua`, `demo/config`): Neovim and the plugins it clones, typing
into `demo/sample/notes.txt`. It needs asciinema, agg and ffmpeg, and writes
`demo/out/`.

## License

MIT, see LICENSE.
