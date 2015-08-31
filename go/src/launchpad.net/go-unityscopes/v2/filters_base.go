package scopes

// Filter is implemented by all scope filter types.
type Filter interface {
	serializeFilter() interface{}
}

type FilterDisplayHints int

const (
	FilterDisplayDefault FilterDisplayHints = 0
	FilterDisplayPrimary FilterDisplayHints = 1 << iota
)

// FilterState represents the current state of a set of filters.
type FilterState map[string]interface{}

type filterBase struct {
	Id           string
	DisplayHints FilterDisplayHints
	FilterType   string
}

type filterWithLabel struct {
	filterBase
	Label string
}

type filterWithOptions struct {
	filterWithLabel
	Options []FilterOption
}

// AddOption adds a new option to the filter.
func (f *filterWithOptions) AddOption(id, label string) {
	f.Options = append(f.Options, FilterOption{
		Id:    id,
		Label: label,
	})
}

func (f *filterWithOptions) isValidOption(optionId interface{}) bool {
	for _, o := range f.Options {
		if o.Id == optionId {
			return true
		}
	}
	return false
}

// HasActiveOption returns true if any of the filters options are active.
func (f *filterWithOptions) HasActiveOption(state FilterState) bool {
	for _, optionId := range f.ActiveOptions(state) {
		if f.isValidOption(optionId) {
			return true
		}
	}
	return false
}

// ActiveOptions returns the filter's active options from the filter state.
func (f *filterWithOptions) ActiveOptions(state FilterState) []string {
	var ret []string
	if state[f.Id] != nil {
		ret = state[f.Id].([]string)
	}
	return ret
}

type FilterOption struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}
