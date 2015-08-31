package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"encoding/json"
	"unsafe"
)

type ActivationStatus int

const (
	ActivationNotHandled ActivationStatus = iota
	ActivationShowDash
	ActivationHideDash
	ActivationShowPreview
	ActivationPerformQuery
)

// ActivationResponse is used as the result of a Activate() or
// PerformAction() call on the scope to instruct the dash on what to
// do next.
type ActivationResponse struct {
	Status    ActivationStatus
	Query     *CannedQuery
	ScopeData interface{}
}

// NewActivationResponse creates an ActivationResponse with the given status
//
// This function should not be used to create an
// ActivationPerformQuery response: use NewActivationResponseForQuery
// instead.
func NewActivationResponse(status ActivationStatus) *ActivationResponse {
	if status == ActivationPerformQuery {
		panic("Use NewActivationResponseFromQuery for PerformQuery responses")
	}
	return &ActivationResponse{
		Status: status,
		Query:  nil,
	}
}

// NewActivationResponseForQuery creates an ActivationResponse that
// performs the given query.
func NewActivationResponseForQuery(query *CannedQuery) *ActivationResponse {
	return &ActivationResponse{
		Status: ActivationPerformQuery,
		Query:  query,
	}
}

func (r *ActivationResponse) update(responsePtr *C._ActivationResponse) error {
	if r.Status == ActivationPerformQuery {
		C.activation_response_init_query(responsePtr, r.Query.q)
	} else {
		C.activation_response_init_status(responsePtr, C.int(r.Status))
	}
	if r.ScopeData != nil {
		data, err := json.Marshal(r.ScopeData)
		if err != nil {
			return err
		}
		var errorString *C.char
		C.activation_response_set_scope_data(responsePtr, (*C.char)(unsafe.Pointer(&data[0])), C.int(len(data)), &errorString)
		if err = checkError(errorString); err != nil {
			return err
		}
	}
	return nil
}

// SetScopeData stores data that will be passed through to the preview
// for ActivationShowPreview type responses.
func (r *ActivationResponse) SetScopeData(v interface{}) {
	r.ScopeData = v
}
