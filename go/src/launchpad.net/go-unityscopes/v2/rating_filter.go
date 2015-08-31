package scopes

import (
//	"sort"
)

// RatingFilter is a filter that allows for rating-based selection
type RatingFilter struct {
	filterWithOptions
	OnIcon  string
	OffIcon string
}

// NewRatingFilter creates a new rating filter.
func NewRatingFilter(id, label string) *RatingFilter {
	return &RatingFilter{
		filterWithOptions: filterWithOptions{
			filterWithLabel: filterWithLabel{
				filterBase: filterBase{
					Id:           id,
					DisplayHints: FilterDisplayDefault,
					FilterType:   "rating",
				},
				Label: label,
			},
		},
	}
}

// ActiveRating gets active option from an instance of FilterState for this filter.
func (f *RatingFilter) ActiveRating(state FilterState) (string, bool) {
	rating, ok := state[f.Id].(string)
	return rating, ok
}

// UpdateState updates the value of a particular option in the filter state.
func (f *RatingFilter) UpdateState(state FilterState, optionId string, active bool) {
	if !f.isValidOption(optionId) {
		panic("invalid option ID")
	}
	// If the state isn't in a form we expect, treat it as empty
	selected, ok := state[f.Id].(string)
	if ok && selected == optionId && active == false {
		delete(state, f.Id)
	} else {
		if active {
			state[f.Id] = optionId
		}
	}
}

func (f *RatingFilter) serializeFilter() interface{} {
	return map[string]interface{}{
		"filter_type":   f.FilterType,
		"id":            f.Id,
		"display_hints": f.DisplayHints,
		"label":         f.Label,
		"options":       f.Options,
	}
}
