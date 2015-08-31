package scopes

import (
	"errors"
	"fmt"
)

type SliderType int

const (
	LessThan SliderType = 0
	MoreThan SliderType = 1 << iota
)

// ValueSliderFilter is a value slider filter that allows for selecting a value within a given range.
type ValueSliderFilter struct {
	filterWithLabel
	Type               SliderType
	DefaultValue       float64
	Min                float64
	Max                float64
	ValueLabelTemplate string
}

// NewValueSliderFilter creates a new value slider filter.
func NewValueSliderFilter(id, label, label_template string, min, max float64) *ValueSliderFilter {
	return &ValueSliderFilter{
		filterWithLabel: filterWithLabel{
			filterBase: filterBase{
				Id:           id,
				DisplayHints: FilterDisplayDefault,
				FilterType:   "value_slider",
			},
			Label: label,
		},
		ValueLabelTemplate: label_template,
		Min:                min,
		Max:                max,
		DefaultValue:       max,
	}
}

// Value gets value of this filter from filter state object.
// If the value is not set for the filter it returns false as the second return statement,
// it returns true otherwise
func (f *ValueSliderFilter) Value(state FilterState) (float64, bool) {
	value, ok := state[f.Id].(float64)
	return value, ok
}

// UpdateState updates the value of the filter to the given value
func (f *ValueSliderFilter) UpdateState(state FilterState, value float64) error {
	if value < f.Min || value > f.Max {
		return errors.New(fmt.Sprintf("ValueSliderFilter:UpdateState: value %f outside of allowed range (%f,%f)", value, f.Min, f.Max))
	}
	state[f.Id] = value
	return nil
}

func (f *ValueSliderFilter) serializeFilter() interface{} {
	return map[string]interface{}{
		"filter_type":    f.FilterType,
		"id":             f.Id,
		"display_hints":  f.DisplayHints,
		"label":          f.Label,
		"label_template": f.ValueLabelTemplate,
		"min":            f.Min,
		"max":            f.Max,
		"default":        f.DefaultValue,
		"slider_type":    f.Type,
	}
}
