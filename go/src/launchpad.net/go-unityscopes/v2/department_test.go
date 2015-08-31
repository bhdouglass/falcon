package scopes_test

import (
	. "gopkg.in/check.v1"
	"launchpad.net/go-unityscopes/v2"
)

func (s *S) TestDepartment(c *C) {
	query := scopes.NewCannedQuery("scope", "query_string", "department_string")
	department, err := scopes.NewDepartment("department_string2", query, "TEST_DEPARTMENT")
	c.Assert(err, IsNil)

	department.SetAlternateLabel("test_alternate_label")
	c.Check(department.AlternateLabel(), Equals, "test_alternate_label")
	c.Check(department.Id(), Equals, "department_string2")
	c.Check(department.Label(), Equals, "TEST_DEPARTMENT")

	department.SetHasSubdepartments(true)
	c.Check(department.HasSubdepartments(), Equals, true)

	department.SetHasSubdepartments(false)
	c.Check(department.HasSubdepartments(), Equals, false)

	department2, err := scopes.NewDepartment("sub_department_string", query, "TEST_SUB_DEPARTMENT")
	c.Assert(err, IsNil)
	department2.SetAlternateLabel("test_alternate_label_2")

	department3, err := scopes.NewDepartment("sub_department_2_string", query, "TEST_SUB_DEPARTMENT_2")
	c.Assert(err, IsNil)
	department3.SetAlternateLabel("test_alternate_label_3")

	subdepartments := department.Subdepartments()
	c.Check(len(subdepartments), Equals, 0)
	c.Check(department.HasSubdepartments(), Equals, false)

	department.SetSubdepartments([]*scopes.Department{department2, department3})
	subdepartments = department.Subdepartments()

	c.Check(len(subdepartments), Equals, 2)
	c.Check(department.HasSubdepartments(), Equals, true)

	// verify that the values are correct in all subdepartments
	c.Check(subdepartments[0].Id(), Equals, department2.Id())
	c.Check(subdepartments[0].Label(), Equals, department2.Label())
	c.Check(subdepartments[0].AlternateLabel(), Equals, department2.AlternateLabel())
	c.Check(subdepartments[1].Id(), Equals, department3.Id())
	c.Check(subdepartments[1].Label(), Equals, department3.Label())
	c.Check(subdepartments[1].AlternateLabel(), Equals, department3.AlternateLabel())

	sub_depts := make([]*scopes.Department, 0)
	department.SetSubdepartments(sub_depts)

	subdepartments = department.Subdepartments()
	c.Check(len(subdepartments), Equals, 0)
	c.Check(department.HasSubdepartments(), Equals, false)

	department.SetSubdepartments([]*scopes.Department{department2, department3})

	subdepartments = department.Subdepartments()
	c.Check(len(subdepartments), Equals, 2)
	c.Check(department.HasSubdepartments(), Equals, true)

	c.Check(subdepartments[0].Id(), Equals, department2.Id())
	c.Check(subdepartments[0].Label(), Equals, department2.Label())
	c.Check(subdepartments[0].AlternateLabel(), Equals, department2.AlternateLabel())
	c.Check(subdepartments[1].Id(), Equals, department3.Id())
	c.Check(subdepartments[1].Label(), Equals, department3.Label())
	c.Check(subdepartments[1].AlternateLabel(), Equals, department3.AlternateLabel())

	stored_query := department.Query()
	c.Check(stored_query.ScopeID(), Equals, "scope")
	c.Check(stored_query.DepartmentID(), Equals, "department_string2")
	c.Check(stored_query.QueryString(), Equals, "query_string")
}

func (s *S) TestDepartmentDifferentCreation(c *C) {
	query := scopes.NewCannedQuery("scope", "query_string", "department_string")
	department, err := scopes.NewDepartment("", query, "TEST_DEPARTMENT")

	c.Assert(err, IsNil)
	c.Check(department.Id(), Equals, "")
	c.Check(department.Label(), Equals, "TEST_DEPARTMENT")
}

func (s *S) TestDepartmentEmptyLabel(c *C) {
	query := scopes.NewCannedQuery("scope", "query_string", "department_string")
	department, err := scopes.NewDepartment("", query, "")
	c.Check(err, Not(Equals), nil)
	c.Check(department, IsNil)

	department, err = scopes.NewDepartment("dept_id", query, "")
	c.Check(err, Not(Equals), nil)
	c.Check(department, IsNil)
}
