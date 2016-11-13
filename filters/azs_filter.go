package filters

type AZsFilter struct {
	azsEnabled map[string]bool
}

func NewAZsFilter(filters []string) *AZsFilter {
	azsEnabled := make(map[string]bool)

	for _, az := range filters {
		azsEnabled[az] = true
	}

	return &AZsFilter{azsEnabled: azsEnabled}
}

func (f *AZsFilter) Enabled(az string) bool {
	if len(f.azsEnabled) == 0 {
		return true
	}

	if f.azsEnabled[az] {
		return true
	}

	return false
}
