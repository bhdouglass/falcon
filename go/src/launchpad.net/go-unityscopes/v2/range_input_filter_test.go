package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestRangeInputFilter(c *C) {
	filter1 := scopes.NewRangeInputFilter("f1", "Options", "start_label", "end_label", "unit_label")
	c.Check("f1", Equals, filter1.Id)
	c.Check("Options", Equals, filter1.Label)
	c.Check(filter1.DisplayHints, Equals, scopes.FilterDisplayDefault)
	c.Check("start_label", Equals, filter1.StartLabel)
	c.Check("end_label", Equals, filter1.EndLabel)
	c.Check("unit_label", Equals, filter1.UnitLabel)

	// check the selection
	fstate := make(scopes.FilterState)
	start, found := filter1.StartValue(fstate)
	c.Check(found, Equals, false)

	end, found := filter1.EndValue(fstate)
	c.Check(found, Equals, false)

	// test setting floats
	err := filter1.UpdateState(fstate, 10.2, 100.4)
	c.Check(err, IsNil)

	start, found = filter1.StartValue(fstate)
	c.Check(start, Equals, 10.2)
	c.Check(found, Equals, true)

	end, found = filter1.EndValue(fstate)
	c.Check(end, Equals, 100.4)
	c.Check(found, Equals, true)

	// test setting floats with no decimals
	err = filter1.UpdateState(fstate, 10.0, 100.0)
	c.Check(err, IsNil)

	start, found = filter1.StartValue(fstate)
	c.Check(start, Equals, 10.0)
	c.Check(found, Equals, true)

	end, found = filter1.EndValue(fstate)
	c.Check(end, Equals, 100.0)
	c.Check(found, Equals, true)

	// test setting mixed floats and integers
	err = filter1.UpdateState(fstate, 10, 100.0)
	c.Check(err, IsNil)

	start, found = filter1.StartValue(fstate)
	c.Check(start, Equals, float64(10))
	c.Check(found, Equals, true)

	end, found = filter1.EndValue(fstate)
	c.Check(end, Equals, 100.0)
	c.Check(found, Equals, true)

	// test integers
	err = filter1.UpdateState(fstate, 10, 100)
	c.Check(err, IsNil)

	start, found = filter1.StartValue(fstate)
	c.Check(start, Equals, float64(10))
	c.Check(found, Equals, true)

	end, found = filter1.EndValue(fstate)
	c.Check(end, Equals, float64(100))
	c.Check(found, Equals, true)

	// test integer and nil
	err = filter1.UpdateState(fstate, nil, 100)
	c.Check(err, IsNil)

	start, found = filter1.StartValue(fstate)
	c.Check(found, Equals, false)

	end, found = filter1.EndValue(fstate)
	c.Check(end, Equals, float64(100))
	c.Check(found, Equals, true)

	// test float and nil
	err = filter1.UpdateState(fstate, 10.4, nil)
	c.Check(err, IsNil)

	start, found = filter1.StartValue(fstate)
	c.Check(start, Equals, 10.4)
	c.Check(found, Equals, true)

	end, found = filter1.EndValue(fstate)
	c.Check(found, Equals, false)

	// test both nil
	err = filter1.UpdateState(fstate, nil, nil)
	c.Check(err, IsNil)

	start, found = filter1.StartValue(fstate)
	c.Check(found, Equals, false)

	end, found = filter1.EndValue(fstate)
	c.Check(found, Equals, false)

	// start greater then end
	err = filter1.UpdateState(fstate, 10, 0.6)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "RangeInputFilter::UpdateState(): start_value 10 is greater or equal to end_value 0.6 for filter f1")

	start, found = filter1.StartValue(fstate)
	c.Check(found, Equals, false)

	end, found = filter1.EndValue(fstate)
	c.Check(found, Equals, false)

	// start equals end
	err = filter1.UpdateState(fstate, 10, 10.0)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "RangeInputFilter::UpdateState(): start_value 10 is greater or equal to end_value 10 for filter f1")

	start, found = filter1.StartValue(fstate)
	c.Check(found, Equals, false)

	end, found = filter1.EndValue(fstate)
	c.Check(found, Equals, false)

	// bad values
	err = filter1.UpdateState(fstate, "", 10.0)
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "RangeInputFilter:UpdateState: Bad type for start value. Valid types are int float64 and nil")

	err = filter1.UpdateState(fstate, 1, "")
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "RangeInputFilter:UpdateState: Bad type for end value. Valid types are int float64 and nil")

	err = filter1.UpdateState(fstate, 1, []int{1, 2})
	c.Check(err, Not(Equals), nil)
	c.Check(err.Error(), Equals, "RangeInputFilter:UpdateState: Bad type for end value. Valid types are int float64 and nil")

	start, found = filter1.StartValue(fstate)
	c.Check(found, Equals, false)

	end, found = filter1.EndValue(fstate)
	c.Check(found, Equals, false)

	// try to set values again
	err = filter1.UpdateState(fstate, 10, 100)
	c.Check(err, IsNil)

	start, found = filter1.StartValue(fstate)
	c.Check(start, Equals, float64(10))
	c.Check(found, Equals, true)

	end, found = filter1.EndValue(fstate)
	c.Check(end, Equals, float64(100))
	c.Check(found, Equals, true)
}
