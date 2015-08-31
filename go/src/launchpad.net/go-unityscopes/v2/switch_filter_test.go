package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestSwitchFilter(c *C) {
	filter1 := scopes.NewSwitchFilter("f1", "Options")
	c.Check("f1", Equals, filter1.Id)
	c.Check("Options", Equals, filter1.Label)
	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayDefault)

	fstate := make(scopes.FilterState)
	c.Check(filter1.IsOn(fstate), Equals, false)

	// set on
	filter1.UpdateState(fstate, true)
	c.Check(filter1.IsOn(fstate), Equals, true)

	filter1.UpdateState(fstate, true)
	filter1.UpdateState(fstate, false)
	filter1.UpdateState(fstate, true)
	c.Check(filter1.IsOn(fstate), Equals, true)

	filter1.UpdateState(fstate, false)
	filter1.UpdateState(fstate, true)
	filter1.UpdateState(fstate, false)
	c.Check(filter1.IsOn(fstate), Equals, false)
}
