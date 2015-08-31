package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestResultSetURI(c *C) {
	r := scopes.NewTestingResult()
	c.Check(r.SetURI("http://example.com"), IsNil)
	c.Check(r.URI(), Equals, "http://example.com")

	var uri string
	c.Check(r.Get("uri", &uri), IsNil)
	c.Check(uri, Equals, "http://example.com")
}

func (s *S) TestResultSetTitle(c *C) {
	r := scopes.NewTestingResult()
	c.Check(r.SetTitle("The title"), IsNil)
	c.Check(r.Title(), Equals, "The title")

	var title string
	c.Check(r.Get("title", &title), IsNil)
	c.Check(title, Equals, "The title")
}

func (s *S) TestResultSetArt(c *C) {
	r := scopes.NewTestingResult()
	c.Check(r.SetArt("http://example.com/foo.png"), IsNil)
	c.Check(r.Art(), Equals, "http://example.com/foo.png")

	var uri string
	c.Check(r.Get("art", &uri), IsNil)
	c.Check(uri, Equals, "http://example.com/foo.png")
}

func (s *S) TestResultSetDndURI(c *C) {
	r := scopes.NewTestingResult()
	c.Check(r.SetDndURI("http://example.com"), IsNil)
	c.Check(r.DndURI(), Equals, "http://example.com")

	var uri string
	c.Check(r.Get("dnd_uri", &uri), IsNil)
	c.Check(uri, Equals, "http://example.com")
}

func (s *S) TestResultSetComplexValue(c *C) {
	type Attr struct {
		Value string `json:"value"`
	}

	r := scopes.NewTestingResult()
	c.Check(r.Set("attributes", []Attr{
		Attr{"one"},
		Attr{"two"},
	}), IsNil)

	// Check that the value has been encoded as expeected:
	var v interface{}
	c.Check(r.Get("attributes", &v), IsNil)
	c.Check(v, DeepEquals, []interface{}{
		map[string]interface{}{"value": "one"},
		map[string]interface{}{"value": "two"},
	})

	// The value can also be decoded into the complex structure too
	var v2 []Attr
	c.Check(r.Get("attributes", &v2), IsNil)
	c.Check(v2, DeepEquals, []Attr{
		Attr{"one"},
		Attr{"two"},
	})
}

func testMarshallingThisFunction(i int) int {
	return i
}

func (s *S) TestResultSetBadValue(c *C) {
	type Attr struct {
		value  int
		value2 float64
	}

	r := scopes.NewTestingResult()
	c.Check(r.Set("attributes", testMarshallingThisFunction), Not(Equals), nil)
}

func (s *S) TestResultGetBadValue(c *C) {
	type Attr struct {
		value  int
		value2 float64
	}

	r := scopes.NewTestingResult()
	var attr string
	c.Check(r.Get("bad_attribute", &attr), Not(Equals), nil)
}
