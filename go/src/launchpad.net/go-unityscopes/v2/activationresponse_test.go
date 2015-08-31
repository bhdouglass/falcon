package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestActivationResponse(c *C) {
	// check all different status
	response := scopes.NewActivationResponse(scopes.ActivationNotHandled)
	c.Check(response.Status, Equals, scopes.ActivationNotHandled)
	c.Check(response.Query, IsNil)
	c.Check(response.ScopeData, IsNil)

	response = scopes.NewActivationResponse(scopes.ActivationShowDash)
	c.Check(response.Status, Equals, scopes.ActivationShowDash)
	c.Check(response.Query, IsNil)
	c.Check(response.ScopeData, IsNil)

	response = scopes.NewActivationResponse(scopes.ActivationShowDash)
	c.Check(response.Status, Equals, scopes.ActivationShowDash)
	c.Check(response.Query, IsNil)
	c.Check(response.ScopeData, IsNil)

	// we should get panic with ActivationPerformQuery
	c.Check(func() { scopes.NewActivationResponse(scopes.ActivationPerformQuery) }, PanicMatches, "Use NewActivationResponseFromQuery for PerformQuery responses")

	// test SetScopeData
	response.SetScopeData("test_string")
	c.Check(response.ScopeData, Equals, "test_string")

	response.SetScopeData(1999)
	c.Check(response.ScopeData, Equals, 1999)

	response.SetScopeData(1.999)
	c.Check(response.ScopeData, Equals, 1.999)

	response.SetScopeData([]string{"test1", "test2"})
	c.Check(response.ScopeData, DeepEquals, []string{"test1", "test2"})

	// test activation response for query
	query := scopes.NewCannedQuery("scope", "query_string", "department_string")
	response_query := scopes.NewActivationResponseForQuery(query)

	c.Check(response_query.Status, Equals, scopes.ActivationPerformQuery)
	c.Check(response_query.Query, Equals, query)
	c.Check(response_query.ScopeData, IsNil)
}
