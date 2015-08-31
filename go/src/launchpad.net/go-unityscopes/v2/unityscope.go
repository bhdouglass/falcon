package scopes

/*
#cgo CXXFLAGS: -std=c++11 -Wall
#cgo pkg-config: libunity-scopes
#include <stdlib.h>
#include "shim.h"
*/
import "C"
import (
	"encoding/json"
	"errors"
	"flag"
	"path"
	"strings"
	"sync"
	"unsafe"
)

func checkError(errorString *C.char) (err error) {
	if errorString != nil {
		err = errors.New(C.GoString(errorString))
		C.free(unsafe.Pointer(errorString))
	}
	return
}

// Category represents a search result category.
type Category struct {
	c C.SharedPtrData
}

func finalizeCategory(cat *Category) {
	C.destroy_category_ptr(&cat.c[0])
}

// Scope defines the interface that scope implementations must implement
type Scope interface {
	SetScopeBase(base *ScopeBase)
	Search(query *CannedQuery, metadata *SearchMetadata, reply *SearchReply, cancelled <-chan bool) error
	Preview(result *Result, metadata *ActionMetadata, reply *PreviewReply, cancelled <-chan bool) error
}

// Activator is an interface that should be implemented by scopes that
// need to handle result activation directly.
type Activator interface {
	Scope
	Activate(result *Result, metadata *ActionMetadata) (*ActivationResponse, error)
}

// PerformActioner is an interface that should be implemented by
// scopes that need to handle preview actions directly.
type PerformActioner interface {
	Scope
	PerformAction(result *Result, metadata *ActionMetadata, widgetId, actionId string) (*ActivationResponse, error)
}

//export callScopeSearch
func callScopeSearch(scope Scope, queryPtr, metadataPtr unsafe.Pointer, replyData *C.uintptr_t, cancel <-chan bool) {
	query := makeCannedQuery((*C._CannedQuery)(queryPtr))
	metadata := makeSearchMetadata((*C._SearchMetadata)(metadataPtr))
	reply := makeSearchReply(replyData)

	go func() {
		err := scope.Search(query, metadata, reply, cancel)
		if err != nil {
			reply.Error(err)
			return
		}
		reply.Finished()
	}()
}

//export callScopePreview
func callScopePreview(scope Scope, resultPtr, metadataPtr unsafe.Pointer, replyData *C.uintptr_t, cancel <-chan bool) {
	result := makeResult((*C._Result)(resultPtr))
	metadata := makeActionMetadata((*C._ActionMetadata)(metadataPtr))
	reply := makePreviewReply(replyData)

	go func() {
		err := scope.Preview(result, metadata, reply, cancel)
		if err != nil {
			reply.Error(err)
			return
		}
		reply.Finished()
	}()
}

//export callScopeActivate
func callScopeActivate(scope Scope, resultPtr, metadataPtr, responsePtr unsafe.Pointer, errorPtr **C.char) {
	switch s := scope.(type) {
	case Activator:
		result := makeResult((*C._Result)(resultPtr))
		metadata := makeActionMetadata((*C._ActionMetadata)(metadataPtr))
		response, err := s.Activate(result, metadata)
		if err == nil {
			err = response.update((*C._ActivationResponse)(responsePtr))
		}
		if err != nil {
			*errorPtr = C.CString(err.Error())
		}
	default:
		// nothing
	}
}

//export callScopePerformAction
func callScopePerformAction(scope Scope, resultPtr, metadataPtr unsafe.Pointer, widgetId, actionId *C.char, responsePtr unsafe.Pointer, errorPtr **C.char) {
	switch s := scope.(type) {
	case PerformActioner:
		result := makeResult((*C._Result)(resultPtr))
		metadata := makeActionMetadata((*C._ActionMetadata)(metadataPtr))
		response, err := s.PerformAction(result, metadata, C.GoString(widgetId), C.GoString(actionId))
		if err == nil {
			err = response.update((*C._ActivationResponse)(responsePtr))
		}
		if err != nil {
			*errorPtr = C.CString(err.Error())
		}
	default:
		// nothing
	}
}

var (
	runtimeConfig = flag.String("runtime", "", "The runtime configuration file for the Unity Scopes library")
	scopeConfig   = flag.String("scope", "", "The scope configuration file for the Unity Scopes library")
)

// ScopeBase exposes information about the scope including settings
// and various directories available for use.
type ScopeBase struct {
	b unsafe.Pointer
}

//export setScopeBase
func setScopeBase(scope Scope, b unsafe.Pointer) {
	if b == nil {
		scope.SetScopeBase(nil)
	} else {
		scope.SetScopeBase(&ScopeBase{b})
	}
}

// ScopeDirectory returns the directory where the scope has been installed
func (b *ScopeBase) ScopeDirectory() string {
	dir := C.scope_base_scope_directory(b.b)
	defer C.free(unsafe.Pointer(dir))
	return C.GoString(dir)
}

// CacheDirectory returns a directory the scope can use to store cache files
func (b *ScopeBase) CacheDirectory() string {
	dir := C.scope_base_cache_directory(b.b)
	defer C.free(unsafe.Pointer(dir))
	return C.GoString(dir)
}

// TmpDirectory returns a directory the scope can use to store temporary files
func (b *ScopeBase) TmpDirectory() string {
	dir := C.scope_base_tmp_directory(b.b)
	defer C.free(unsafe.Pointer(dir))
	return C.GoString(dir)
}

// ListRegistryScopes lists all the scopes existing in the registry
func (b *ScopeBase) ListRegistryScopes() map[string]*ScopeMetadata {
	var nb_scopes C.int
	var c_array **C._ScopeMetadata = C.list_registry_scopes_metadata(b.b, &nb_scopes)
	defer C.free(unsafe.Pointer(c_array))

	length := int(nb_scopes)
	// create a very big slice and then slice it to the number of scopes metadata
	slice := (*[1 << 27]*C._ScopeMetadata)(unsafe.Pointer(c_array))[:length:length]
	var scopesList = make(map[string]*ScopeMetadata)

	for i := 0; i < length; i++ {
		json_data := C.get_scope_metadata_serialized(slice[i])
		defer C.free(unsafe.Pointer(json_data))

		metadata := makeScopeMetadata(slice[i], C.GoString(json_data))
		scopesList[metadata.ScopeId] = metadata
	}

	return scopesList
}

// Settings returns the scope's settings.  The settings will be
// decoded into the given value according to the same rules used by
// json.Unmarshal().
func (b *ScopeBase) Settings(value interface{}) error {
	var length C.int
	data := C.scope_base_settings(b.b, &length)
	defer C.free(data)
	return json.Unmarshal(C.GoBytes(data, length), value)
}

/*
Run will initialise the scope runtime and make a scope availble.  It
is intended to be called from the program's main function, and will
run until the program exits.
*/
func Run(scope Scope) error {
	if !flag.Parsed() {
		flag.Parse()
	}
	if *scopeConfig == "" {
		return errors.New("Scope configuration file not set on command line")
	}
	base := path.Base(*scopeConfig)
	if !strings.HasSuffix(base, ".ini") {
		return errors.New("Scope configuration file does not end in '.ini'")
	}
	scopeId := base[:len(base)-len(".ini")]

	var errorString *C.char
	C.run_scope(unsafe.Pointer(&scopeId), unsafe.Pointer(runtimeConfig), unsafe.Pointer(scopeConfig), unsafe.Pointer(&scope), &errorString)
	return checkError(errorString)
}

var (
	cancelChannels     = make(map[chan bool]bool)
	cancelChannelsLock sync.Mutex
)

//export makeCancelChannel
func makeCancelChannel() chan bool {
	ch := make(chan bool, 1)
	cancelChannelsLock.Lock()
	cancelChannels[ch] = true
	cancelChannelsLock.Unlock()
	return ch
}

//export sendCancelChannel
func sendCancelChannel(ch chan bool) {
	ch <- true
}

//export releaseCancelChannel
func releaseCancelChannel(ch chan bool) {
	cancelChannelsLock.Lock()
	delete(cancelChannels, ch)
	cancelChannelsLock.Unlock()
}
