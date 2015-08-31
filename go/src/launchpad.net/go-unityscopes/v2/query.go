package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"encoding/json"
	"runtime"
	"unsafe"
)

// CannedQuery represents a search query from the user.
type CannedQuery struct {
	q *C._CannedQuery
}

func finalizeCannedQuery(query *CannedQuery) {
	if query.q != nil {
		C.destroy_canned_query(query.q)
	}
	query.q = nil
}

func makeCannedQuery(q *C._CannedQuery) *CannedQuery {
	query := new(CannedQuery)
	runtime.SetFinalizer(query, finalizeCannedQuery)
	query.q = q
	return query
}

// NewCannedQuery creates a new CannedQuery with the given scope ID,
// query string and department ID.
func NewCannedQuery(scopeID, queryString, departmentID string) *CannedQuery {
	return makeCannedQuery(C.new_canned_query(
		unsafe.Pointer(&scopeID),
		unsafe.Pointer(&queryString),
		unsafe.Pointer(&departmentID)))
}

// ScopeID returns the scope ID for this canned query.
func (query *CannedQuery) ScopeID() string {
	s := C.canned_query_get_scope_id(query.q)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// DepartmentID returns the department ID for this canned query.
func (query *CannedQuery) DepartmentID() string {
	s := C.canned_query_get_department_id(query.q)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// QueryString returns the query string for this canned query.
func (query *CannedQuery) QueryString() string {
	s := C.canned_query_get_query_string(query.q)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

// FilterState returns the state of the filters for this canned query.
func (query *CannedQuery) FilterState() FilterState {
	var length C.int
	s := C.canned_query_get_filter_state(query.q, &length)
	defer C.free(s)
	var state FilterState
	if err := json.Unmarshal(C.GoBytes(s, length), &state); err != nil {
		panic(err)
	}
	return state
}

// SetDepartmentID changes the department ID for this canned query.
func (query *CannedQuery) SetDepartmentID(departmentID string) {
	C.canned_query_set_department_id(query.q, unsafe.Pointer(&departmentID))
}

// SetQueryString changes the query string for this canned query.
func (query *CannedQuery) SetQueryString(queryString string) {
	C.canned_query_set_query_string(query.q, unsafe.Pointer(&queryString))
}

// ToURI formats the canned query as a URI.
func (query *CannedQuery) ToURI() string {
	s := C.canned_query_to_uri(query.q)
	defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}
