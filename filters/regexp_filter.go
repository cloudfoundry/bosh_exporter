package filters

import (
	"regexp"
)

type RegexpFilter struct {
	reFilters []*regexp.Regexp
}

func NewRegexpFilter(filters []string) (*RegexpFilter, error) {
	reFilters := []*regexp.Regexp{}

	for _, filter := range filters {
		re, err := regexp.Compile(filter)
		if err != nil {
			return nil, err
		}
		reFilters = append(reFilters, re)
	}

	return &RegexpFilter{reFilters: reFilters}, nil
}

func (f *RegexpFilter) Enabled(expr string) bool {
	if len(f.reFilters) == 0 {
		return true
	}

	for _, re := range f.reFilters {
		matched := re.MatchString(expr)
		if matched {
			return true
		}
	}

	return false
}
