// Command demo is the goda-cmp configuration used by demo/record.sh: the path
// source, narrowed down by the fuzzy filters of goda-cmp-ext-fuzzy.vim and
// shown with the UI pum. The converter mtime shows the modification time
// that the source keeps in path.Data.
package main

import (
	"context"

	"github.com/Shougo/goda-cmp-ui-pum.vim/ext/ui/pum"
	"github.com/Shougo/goda-cmp.vim/cmp"
	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/matcherfuzzy"
	"github.com/delphinus/goda-cmp-ext-fuzzy.vim/ext/filter/sorterfuzzy"

	"github.com/delphinus/goda-cmp-ext-path.vim/ext/source/path"
)

// mtime puts the modification time of a path in the menu of its candidate.
type mtime struct{}

func (mtime) Name() string { return "mtime" }

func (mtime) Filter(_ context.Context, in cmp.FilterInput) ([]cmp.Item, error) {
	out := make([]cmp.Item, len(in.Items))
	for i, it := range in.Items {
		if d, ok := it.UserData.(path.Data); ok {
			it.Menu = d.MTime.Format("2006-01-02 15:04")
		}
		out[i] = it
	}
	return out, nil
}

func main() {
	b := cmp.New()
	b.AddUI(pum.New())
	b.AddSource(path.New())
	b.AddFilter(matcherfuzzy.New())
	b.AddFilter(sorterfuzzy.New())
	b.AddFilter(mtime{})
	b.PatchGlobal(cmp.Config{
		UI:      "pum",
		Sources: []string{"path"},
		SourceOptions: map[string]cmp.SourceOption{
			"path": {
				Mark:                   "P",
				Matchers:               []string{"matcherfuzzy"},
				Sorters:                []string{"sorterfuzzy"},
				Converters:             []string{"mtime"},
				ForceCompletionPattern: path.ForceCompletionPattern,
			},
		},
	})
	b.Main()
}
