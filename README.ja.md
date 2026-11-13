# goda-cmp-ext-path.vim

[English](README.md)

> [!NOTE]
> goda-cmp.vim 自体が開発の途中のため、それに合わせてこのプラグインの API も変わることがあります。

https://github.com/user-attachments/assets/9c510911-74ec-4d9f-b9dc-0944b1c3951b

*名前全体を更新日時と一緒に出し、[goda-cmp-ext-fuzzy.vim](https://github.com/delphinus/goda-cmp-ext-fuzzy.vim) の `matcherfuzzy` で絞り込んだところ。ドットファイルはドットを打つと出る。menu の日時は、デモの converter が `path.Data` から作っている (`demo/record.sh` で録画)。*

[goda-cmp.vim](https://github.com/Shougo/goda-cmp.vim) 用のパスのソース。

ディレクトリの中の名前をそのまま候補にし、フルパス、更新日時、ディレクトリかどうかを各候補の UserData (`path.Data`) に置く。UI はそれを読んで表示できる。

- 補完するのは、カーソルの前の語が `/`、`~/`、`./`、`../`、`$VAR/`、`${VAR}/` で始まるときだけ。`github.com/foo` のような語はパスとみなさない。
- 相対パスはバッファのあるディレクトリを起点にする。
- ドットファイルは、ドットを打ったときだけ出す。

## goda-cmp-ext の `file` ソースとの違い

[goda-cmp-ext](https://github.com/Shougo/goda-cmp-ext) の `file` ソースは、打った部分を除いた名前の残りだけを候補にし、候補にフルパスを持たせない。このソースは名前全体を候補にするので、fuzzy なマッチャで絞り込める。UI は `path.Data` からフルパスと更新日時を読める。

## インストール

プラグインマネージャで goda-cmp.vim と一緒に入れる。`g:goda#workspace` があれば (goda-pm.vim が設定する)、`goda#build()` は入れたプラグインで設定をビルドする。そうでなければ、設定の Go のモジュールに足す:

```sh
go get github.com/delphinus/goda-cmp-ext-path.vim
```

## 設定

設定の `main.go` で:

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
		// パスを打っているあいだは path の候補だけを出す。
		DynamicSources: func(_ context.Context, _ sdk.Host, a cmp.DynamicSourcesArgs) ([]string, error) {
			if a.Context.Mode == "i" && path.InContext(a.Context.Input) {
				return []string{"path"}, nil
			}
			return nil, nil
		},
		SourceOptions: map[string]cmp.SourceOption{
			// / を打った直後から補完する。
			"path": {ForceCompletionPattern: path.ForceCompletionPattern},
		},
	})
	b.Main()
}
```

## デモ

`demo/record.sh` は、最小の設定 (`demo/init.lua`、`demo/config`) でデモを録画する。Neovim を使い、プラグインは自分で取得して、`demo/sample/notes.txt` に打ち込む。asciinema、agg、ffmpeg が要る。出力は `demo/out/`。

## ライセンス

MIT。LICENSE を参照。
