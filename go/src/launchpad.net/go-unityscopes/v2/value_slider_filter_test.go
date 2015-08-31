package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestValueSliderFilter(c *C) {
	filter1 := scopes.NewValueSliderFilter("f1", "Options", "label", 10.0, 100.0)
	c.Check("f1", Equals, filter1.Id)
	c.Check("Options", Equals, filter1.Label)
	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayDefault)
	c.Check(filter1.DefaultValue, Equals, 100.0)
	c.Check(filter1.Min, Equals, 10.0)
	c.Check(filter1.Max, Equals, 100.0)
	c.Check(filter1.ValueLabelTemplate, Equals, "label")

	fstate := make(scopes.FilterState)
	value, ok := filter1.Value(fstate)
	c.Check(value, Equals, 0.0)
	c.Check(ok, Equals, false)

	err := filter1.UpdateState(fstate, 30.5)
	c.Check(err, IsNil)
	value, ok = filter1.Value(fstate)
	c.Check(value, Equals, 30.5)
	c.Check(ok, Equals, true)

	err = filter1.UpdateState(fstate, 44.5)
	c.Check(err, IsNil)
	value, ok = filter1.Value(fstate)
	c.Check(value, Equals, 44.5)
	c.Check(ok, Equals, true)

	err = filter1.UpdateState(fstate, 3545.33)
	c.Check(err, Not(Equals), nil)
	value, ok = filter1.Value(fstate)
	c.Check(value, Equals, 44.5)
	c.Check(ok, Equals, true)
}
