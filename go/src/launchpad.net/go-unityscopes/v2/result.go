package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"encoding/json"
	"runtime"
	"unsafe"
)

// Result represents a result from the scope
type Result struct {
	result *C._Result
}

func makeResult(res *C._Result) *Result {
	result := new(Result)
	runtime.SetFinalizer(result, finalizeResult)
	result.result = res
	return result
}

func finalizeResult(res *Result) {
	if res.result != nil {
		C.destroy_result(res.result)
	}
	res.result = nil
}

// Get returns the named result attribute.
//
// The value is decoded into the variable pointed to by the second
// argument.  If the types do not match, an error will be returned.
//
// If the attribute does not exist, an error is returned.
func (res *Result) Get(attr string, value interface{}) error {
	var (
		length      C.int
		errorString *C.char
	)
	data := C.result_get_attr(res.result, unsafe.Pointer(&attr), &length, &errorString)
	if err := checkError(errorString); err != nil {
		return err
	}
	defer C.free(data)
	return json.Unmarshal(C.GoBytes(data, length), value)
}

// Set sets the named result attribute.
//
// An error may be returned if the value can not be stored, or if
// there is any other problems updating the result.
func (res *Result) Set(attr string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	stringValue := string(data)

	var errorString *C.char
	C.result_set_attr(res.result, unsafe.Pointer(&attr), unsafe.Pointer(&stringValue), &errorString)
	return checkError(errorString)
}

// SetInterceptActivation marks this result as needing custom activation handling.
//
// By default, results are activated by the client directly (e.g. by
// running the application associated with the result URI).  For
// results with this flag set though, the scope will be asked to
// perform activation and should implement the Activate method.
func (res *Result) SetInterceptActivation() {
	C.result_set_intercept_activation(res.result)
}

// SetURI sets the "uri" attribute of the result.
func (res *Result) SetURI(uri string) error {
	return res.Set("uri", uri)
}

// SetTitle sets the "title" attribute of the result.
func (res *Result) SetTitle(title string) error {
	return res.Set("title", title)
}

// SetArt sets the "art" attribute of the result.
func (res *Result) SetArt(art string) error {
	return res.Set("art", art)
}

// SetDndURI sets the "dnd_uri" attribute of the result.
func (res *Result) SetDndURI(uri string) error {
	return res.Set("dnd_uri", uri)
}

func (res *Result) getString(attr string) string {
	var value string
	if err := res.Get(attr, &value); err != nil {
		return ""
	}
	return value
}

// URI returns the "uri" attribute of the result if set, or an empty string.
func (res *Result) URI() string {
	return res.getString("uri")
}

// Title returns the "title" attribute of the result if set, or an empty string.
func (res *Result) Title() string {
	return res.getString("title")
}

// Art returns the "art" attribute of the result if set, or an empty string.
func (res *Result) Art() string {
	return res.getString("art")
}

// DndURI returns the "dnd_uri" attribute of the result if set, or an
// empty string.
func (res *Result) DndURI() string {
	return res.getString("dnd_uri")
}

// CategorisedResult represents a result linked to a particular category.
//
// CategorisedResult embeds Result, so all of its attribute
// manipulation methods can be used on variables of this type.
type CategorisedResult struct {
	Result
}

// NewCategorisedResult creates a new empty result linked to the given
// category.
func NewCategorisedResult(category *Category) *CategorisedResult {
	res := new(CategorisedResult)
	runtime.SetFinalizer(res, finalizeCategorisedResult)
	res.result = C.new_categorised_result(&category.c[0])
	return res
}

func finalizeCategorisedResult(res *CategorisedResult) {
	finalizeResult(&res.Result)
}
