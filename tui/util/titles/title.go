package titles

import (
	"github.com/minkezhang/truffle-api/db/atom"
	"sort"
)

var (
	priority = map[string]int{
		"en": 0,
		"":   1,
		"ja": 2,
	}
)

func Sort(titles []atom.T) []atom.T {
	titles = append([]atom.T{}, titles...)
	f := func(i, j int) bool {
		if titles[i].Localization == titles[j].Localization {
			return titles[i].Title < titles[j].Title
		}
		u, ok := priority[titles[i].Localization]
		if !ok {
			u = int(^uint(0) >> 1)
		}
		v, ok := priority[titles[j].Localization]
		if !ok {
			v = int(^uint(0) >> 1)
		}
		return u < v
	}
	sort.SliceStable(titles, f)
	return titles
}

func Title(titles []atom.T) string { return Sort(titles)[0].Title }
