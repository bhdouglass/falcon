package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestRadioButtonsFilter(c *C) {
	filter1 := scopes.NewRadioButtonsFilter("f1", "Options")
	c.Check("f1", Equals, filter1.Id)
	c.Check("Options", Equals, filter1.Label)
	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayDefault)

	filter1.DisplayHints = scopes.FilterDisplayPrimary
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")
	filter1.AddOption("3", "Option 3")

	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayPrimary)

	// verify the list of options
	c.Check(len(filter1.Options), Equals, 3)
	c.Check(filter1.Options, DeepEquals, []scopes.FilterOption{scopes.FilterOption{"1", "Option 1"},
		scopes.FilterOption{"2", "Option 2"},
		scopes.FilterOption{"3", "Option 3"}})

	// check the selection
	fstate := make(scopes.FilterState)
	c.Check(filter1.HasActiveOption(fstate), Equals, false)
}

func (s *S) TestRadioButtonsFilterSingleSelection(c *C) {
	filter1 := scopes.NewRadioButtonsFilter("f1", "Options")
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")
	filter1.AddOption("3", "Option 2")

	fstate := make(scopes.FilterState)
	_, ok := fstate["route"]
	c.Check(ok, Equals, false)
	c.Check(filter1.HasActiveOption(fstate), Equals, false)

	// enable option1
	filter1.UpdateState(fstate, "1", true)
	_, ok = fstate["f1"]
	c.Check(ok, Equals, true)
	c.Check(filter1.HasActiveOption(fstate), Equals, true)

	active := filter1.ActiveOptions(fstate)
	c.Check(len(active), Equals, 1)
	c.Check(active, DeepEquals, []string{"1"})

	// enable option2, option1 get disabled
	filter1.UpdateState(fstate, "2", true)
	active = filter1.ActiveOptions(fstate)
	c.Check(len(active), Equals, 1)
	c.Check(active, DeepEquals, []string{"2"})

	// disable option1; filter state remains in the FilterState, just no options are selected
	filter1.UpdateState(fstate, "2", false)
	_, ok = fstate["f1"]
	c.Check(ok, Equals, true)
	active = filter1.ActiveOptions(fstate)
	c.Check(len(active), Equals, 0)
}

func (s *S) TestRadioButtonsFilterBadOption(c *C) {
	filter1 := scopes.NewRadioButtonsFilter("f1", "Options")
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")
	filter1.AddOption("3", "Option 3")

	fstate := make(scopes.FilterState)

	c.Assert(func() { filter1.UpdateState(fstate, "5", true) }, PanicMatches, "invalid option ID")
	c.Assert(func() { filter1.UpdateState(fstate, "5", false) }, PanicMatches, "invalid option ID")
}
