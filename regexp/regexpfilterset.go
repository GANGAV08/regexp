package regexp

import (
	"regexp"

	"github.com/golang/groupcache/lru"

	"fmt"
)

type FilterSet struct {
	regexes      []*regexp.Regexp
	cacheEnabled bool
	cache        *lru.Cache
}

func SayHello() {
	fmt.Println("Hello Go!")
}

func NewFilterSet(filters []string, cfg *Config) (*FilterSet, error) {
	fs := &FilterSet{
		regexes: make([]*regexp.Regexp, 0, len(filters)),
	}

	if cfg != nil && cfg.CacheEnabled {
		fs.cacheEnabled = true
		fs.cache = lru.New(cfg.CacheMaxNumEntries)
	}

	if err := fs.addFilters(filters); err != nil {
		return nil, err
	}

	return fs, nil
}

func (rfs *FilterSet) Matches(toMatch string) bool {
	if rfs.cacheEnabled {
		if v, ok := rfs.cache.Get(toMatch); ok {
			return v.(bool)
		}
	}

	for _, r := range rfs.regexes {
		if r.MatchString(toMatch) {
			if rfs.cacheEnabled {
				rfs.cache.Add(toMatch, true)
			}
			return true
		}
	}

	if rfs.cacheEnabled {
		rfs.cache.Add(toMatch, false)
	}
	return false
}

// addFilters compiles all the given filters and stores them as regexes.
// All regexes are automatically anchored to enforce full string matches.
func (rfs *FilterSet) addFilters(filters []string) error {
	dedup := make(map[string]struct{}, len(filters))
	for _, f := range filters {
		if _, ok := dedup[f]; ok {
			continue
		}

		re, err := regexp.Compile(f)
		if err != nil {
			return err
		}
		rfs.regexes = append(rfs.regexes, re)
		dedup[f] = struct{}{}
	}

	return nil
}
