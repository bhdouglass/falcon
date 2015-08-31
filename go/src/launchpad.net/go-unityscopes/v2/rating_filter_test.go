package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestRatingFilter(c *C) {
	filter1 := scopes.NewRatingFilter("f1", "Options")
	c.Check("f1", Equals, filter1.Id)
	c.Check("Options", Equals, filter1.Label)
	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayDefault)

	filter1.DisplayHints = scopes.FilterDisplayPrimary
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")

	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayPrimary)
	c.Check(2, Equals, len(filter1.Options))
	c.Check("1", Equals, filter1.Options[0].Id)
	c.Check("Option 1", Equals, filter1.Options[0].Label)
	c.Check("2", Equals, filter1.Options[1].Id)
	c.Check("Option 2", Equals, filter1.Options[1].Label)

	// verify the list of options
	c.Check(len(filter1.Options), Equals, 2)
	c.Check(filter1.Options, DeepEquals, []scopes.FilterOption{scopes.FilterOption{"1", "Option 1"}, scopes.FilterOption{"2", "Option 2"}})
}

func (s *S) TestRatingFilterMultiSelection(c *C) {
	filter1 := scopes.NewRatingFilter("f1", "Options")
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")
	filter1.AddOption("3", "Option 3")

	fstate := make(scopes.FilterState)

	// enable option1 & option2
	filter1.UpdateState(fstate, "1", true)
	_, ok := fstate["f1"]
	c.Check(ok, Equals, true)

	active, ok := filter1.ActiveRating(fstate)
	c.Check(active, Equals, "1")
	c.Check(ok, Equals, true)

	// disable option1
	filter1.UpdateState(fstate, "1", false)
	active, ok = filter1.ActiveRating(fstate)
	c.Check(active, Equals, "")
	c.Check(ok, Equals, false)

	filter1.UpdateState(fstate, "3", true)

	active, ok = filter1.ActiveRating(fstate)
	c.Check(active, Equals, "3")
	c.Check(ok, Equals, true)

	// select another one
	filter1.UpdateState(fstate, "1", true)
	active, ok = filter1.ActiveRating(fstate)
	c.Check(active, Equals, "1")
	c.Check(ok, Equals, true)

	// select another one
	filter1.UpdateState(fstate, "2", true)

	// erase not selected
	filter1.UpdateState(fstate, "1", false)
	active, ok = filter1.ActiveRating(fstate)
	c.Check(active, Equals, "2")
	c.Check(ok, Equals, true)

	// erase the active one
	filter1.UpdateState(fstate, "2", false)
	active, ok = filter1.ActiveRating(fstate)
	c.Check(ok, Equals, false)
	c.Check(active, Equals, "")
}

func (s *S) TestRatingFilterBadOption(c *C) {
	filter1 := scopes.NewRatingFilter("f1", "Options")
	filter1.AddOption("1", "Option 1")
	filter1.AddOption("2", "Option 2")
	filter1.AddOption("3", "Option 3")

	fstate := make(scopes.FilterState)

	c.Assert(func() { filter1.UpdateState(fstate, "5", true) }, PanicMatches, "invalid option ID")
	c.Assert(func() { filter1.UpdateState(fstate, "5", false) }, PanicMatches, "invalid option ID")
	c.Assert(func() { filter1.UpdateState(fstate, "", false) }, PanicMatches, "invalid option ID")
}
