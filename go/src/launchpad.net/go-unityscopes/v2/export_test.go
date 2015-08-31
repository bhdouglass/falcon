package scopes

import (
	"encoding/json"
)

// This file exports certain private functions for use by tests.

func NewTestingResult() *Result {
	return newTestingResult()
}

func NewTestingScopeMetadata(json_data string) ScopeMetadata {
	var scopeMetadata ScopeMetadata
	if err := json.Unmarshal([]byte(json_data), &scopeMetadata); err != nil {
		panic(err)
	}

	return scopeMetadata
}
