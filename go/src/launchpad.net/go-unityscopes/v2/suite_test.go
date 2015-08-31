package scopes_test

import (
	"errors"
	. "gopkg.in/check.v1"
	"testing"
)

// The following type is used only for testing unserializable cases
type unserializable struct{}

func (u unserializable) MarshalJSON() ([]byte, error) {
	return nil, errors.New("Can not marshal to JSON")
}

func (u *unserializable) UnmarshalJSON(data []byte) error {
	return errors.New("Can not unmarshal from JSON")
}

type S struct{}

func init() {
	Suite(&S{})
}

func TestAll(t *testing.T) {
	TestingT(t)
}
