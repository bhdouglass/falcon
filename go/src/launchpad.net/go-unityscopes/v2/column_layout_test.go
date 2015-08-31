package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestColumnLayout(c *C) {
	layout := scopes.NewColumnLayout(3)

	c.Check(layout.Size(), Equals, 0)
	c.Check(layout.NumberOfColumns(), Equals, 3)

	err := layout.AddColumn("widget_1", "widget_2")
	c.Assert(err, IsNil)

	c.Check(layout.Size(), Equals, 1)
	c.Check(layout.NumberOfColumns(), Equals, 3)

	col, err := layout.Column(0)
	c.Assert(err, IsNil)

	c.Check(len(col), Equals, 2)
	c.Check(col, DeepEquals, []string{"widget_1", "widget_2"})

	// add another column
	err = layout.AddColumn("widget_3", "widget_4", "widget_5")
	c.Assert(err, IsNil)

	col, err = layout.Column(1)
	c.Assert(err, IsNil)

	c.Check(len(col), Equals, 3)
	c.Check(col, DeepEquals, []string{"widget_3", "widget_4", "widget_5"})

	// check for a bad column
	_, err = layout.Column(2)
	c.Assert(err, Not(Equals), nil)

	// now add the last column
	err = layout.AddColumn("widget_6")
	c.Assert(err, IsNil)

	col, err = layout.Column(2)
	c.Assert(err, IsNil)

	c.Check(len(col), Equals, 1)
	c.Check(col[0], Equals, "widget_6")

	// try to add more columns ... should obtain an error
	err = layout.AddColumn("widget_3", "widget_4", "widget_5")
	c.Assert(err, Not(Equals), nil)

	// check size again
	c.Check(layout.Size(), Equals, 3)
	c.Check(layout.NumberOfColumns(), Equals, 3)

	// check empty list
	layout1col := scopes.NewColumnLayout(1)
	err = layout1col.AddColumn()
	c.Check(err, IsNil)
}
