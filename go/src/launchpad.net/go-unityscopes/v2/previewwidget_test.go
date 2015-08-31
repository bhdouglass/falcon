package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestPreviewWidget(c *C) {
	widget := scopes.NewPreviewWidget("widget_id", "widget_type")
	c.Check(widget.Id(), Equals, "widget_id")
	c.Check(widget.WidgetType(), Equals, "widget_type")

	_, ok := widget["attr_id"]
	c.Check(ok, Equals, false)
	widget.AddAttributeValue("attr_id", "attr_value")

	value, ok := widget["attr_id"]
	c.Check(ok, Equals, true)
	c.Check(value, Equals, "attr_value")

	// set complex value
	widget.AddAttributeValue("attr_slice", []string{"test1", "test2", "test3"})
	value, ok = widget["attr_slice"]
	c.Check(ok, Equals, true)
	c.Check(value, DeepEquals, []string{"test1", "test2", "test3"})

	// attribute mapping

	_, ok = widget["components"]
	c.Check(ok, Equals, false)
	widget.AddAttributeMapping("map_key", "mapping_value")

	value, ok = widget["components"]
	c.Check(ok, Equals, true)

	components := value.(map[string]interface{})

	_, ok = components["map_key_error"]
	c.Check(ok, Equals, false)

	value, ok = components["map_key"]
	c.Check(ok, Equals, true)
	c.Check(value, Equals, "mapping_value")

	// add nother mapping
	widget.AddAttributeMapping("map_key_2", "mapping_value_2")
	value, ok = widget["components"]
	c.Check(ok, Equals, true)
	components = value.(map[string]interface{})

	value, ok = components["map_key"]
	c.Check(ok, Equals, true)
	c.Check(value, Equals, "mapping_value")

	value, ok = components["map_key_2"]
	c.Check(ok, Equals, true)
	c.Check(value, Equals, "mapping_value_2")
}

func (s *S) TestPreviewWidgetAddWidgets(c *C) {
	widget := scopes.NewPreviewWidget("widget_id", "widget_type")
	c.Check(widget.Id(), Equals, "widget_id")
	c.Check(widget.WidgetType(), Equals, "widget_type")

	sub_widget_1 := scopes.NewPreviewWidget("widget1", "expandable")

	// check panic error when adding widget to non expandable
	c.Assert(func() { widget.AddWidget(sub_widget_1) }, PanicMatches, "Can only add widgets to expandable type widgets")

	// check it does not have subwidgets
	_, ok := sub_widget_1["widgets"]
	c.Check(ok, Equals, false)

	sub_widget_11 := scopes.NewPreviewWidget("widget11", "audio")
	sub_widget_12 := scopes.NewPreviewWidget("widget12", "video")

	sub_widget_1.AddWidget(sub_widget_11)
	// now it does have widgets
	widgets, ok := sub_widget_1["widgets"]
	c.Check(ok, Equals, true)
	c.Check(widgets, DeepEquals, []scopes.PreviewWidget{sub_widget_11})

	sub_widget_1.AddWidget(sub_widget_12)
	widgets, ok = sub_widget_1["widgets"]
	c.Check(ok, Equals, true)
	c.Check(widgets, DeepEquals, []scopes.PreviewWidget{sub_widget_11, sub_widget_12})

	main_widget := scopes.NewPreviewWidget("main_widget_id", "expandable")
	// check it does not have subwidgets
	_, ok = main_widget["widgets"]
	c.Check(ok, Equals, false)

	main_widget.AddWidget(sub_widget_1)
	widgets, ok = main_widget["widgets"]
	c.Check(ok, Equals, true)
	c.Check(widgets, DeepEquals, []scopes.PreviewWidget{sub_widget_1})
}
