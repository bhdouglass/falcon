package scopes

import (
	"sort"
)

// OptionSelectorFilter is used to implement single-select or multi-select filters.
type OptionSelectorFilter struct {
	filterWithOptions
	MultiSelect bool
}

// NewOptionSelectorFilter creates a new option filter.
func NewOptionSelectorFilter(id, label string, multiSelect bool) *OptionSelectorFilter {
	return &OptionSelectorFilter{
		filterWithOptions: filterWithOptions{
			filterWithLabel: filterWithLabel{
				filterBase: filterBase{
					Id:           id,
					DisplayHints: FilterDisplayDefault,
					FilterType:   "option_selector",
				},
				Label: label,
			},
		},
		MultiSelect: multiSelect,
	}
}

// UpdateState updates the value of a particular option in the filter state.
func (f *OptionSelectorFilter) UpdateState(state FilterState, optionId string, active bool) {
	if !f.isValidOption(optionId) {
		panic("invalid option ID")
	}
	// For single-select filters, clear the previous state when
	// setting a new active option.
	if active && !f.MultiSelect {
		delete(state, f.Id)
	}
	// If the state isn't in a form we expect, treat it as empty
	selected, _ := state[f.Id].([]string)
	sort.Strings(selected)
	pos := sort.SearchStrings(selected, optionId)
	if active {
		if pos == len(selected) {
			selected = append(selected, optionId)
		} else if pos < len(selected) && selected[pos] != optionId {
			selected = append(selected[:pos], append([]string{optionId}, selected[pos:]...)...)
		}
	} else {
		if pos < len(selected) {
			selected = append(selected[:pos], selected[pos+1:]...)
		}
	}
	state[f.Id] = selected
}

func (f *OptionSelectorFilter) serializeFilter() interface{} {
	return map[string]interface{}{
		"filter_type":   f.FilterType,
		"id":            f.Id,
		"display_hints": f.DisplayHints,
		"label":         f.Label,
		"multi_select":  f.MultiSelect,
		"options":       f.Options,
	}
}
