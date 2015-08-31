package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestQuery(c *C) {
	query := scopes.NewCannedQuery("scope", "query_string", "department_string")

	// basic check
	c.Check(query.ScopeID(), Equals, "scope")
	c.Check(query.DepartmentID(), Equals, "department_string")
	c.Check(query.QueryString(), Equals, "query_string")

	// verify uri
	c.Check(query.ToURI(), Equals, "scope://scope?q=query%5Fstring&dep=department%5Fstring")

	// check setters
	query.SetDepartmentID("department_id")
	c.Check(query.DepartmentID(), Equals, "department_id")

	query.SetQueryString("new_query_value")
	c.Check(query.QueryString(), Equals, "new_query_value")

	// TODO FilterState setter
}
