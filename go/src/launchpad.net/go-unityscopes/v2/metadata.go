package scopes

// #include <stdlib.h>
// #include "shim.h"
import "C"
import (
	"encoding/json"
	"fmt"
	"runtime"
	"unsafe"
)

type ConnectivityStatus int

const (
	ConnectivityStatusUnknown      ConnectivityStatus = 0
	ConnectivityStatusConnected    ConnectivityStatus = 1
	ConnectivityStatusDisconnected ConnectivityStatus = 2
)

// queryMetadata is the base class for extra metadata passed to scopes as a part of a request.
// This base class is not exported
type queryMetadata struct {
	m *C._QueryMetadata
}

// Locale returns the expected locale for the search request.
func (metadata *queryMetadata) Locale() string {
	locale := C.query_metadata_get_locale(metadata.m)
	defer C.free(unsafe.Pointer(locale))
	return C.GoString(locale)
}

// FormFactor returns the form factor for the search request.
func (metadata *queryMetadata) FormFactor() string {
	formFactor := C.query_metadata_get_form_factor(metadata.m)
	defer C.free(unsafe.Pointer(formFactor))
	return C.GoString(formFactor)
}

// SetInternetConnectivity indicates the internet connectivity status.
func (metadata *queryMetadata) SetInternetConnectivity(status ConnectivityStatus) {
	C.query_metadata_set_internet_connectivity(metadata.m, C.int(status))
}

// InternetConnectivity gets internet connectivity status.
func (metadata *queryMetadata) InternetConnectivity() ConnectivityStatus {
	return ConnectivityStatus(C.query_metadata_get_internet_connectivity(metadata.m))
}

// SearchMetadata holds additional metadata about the search request.
type SearchMetadata struct {
	queryMetadata
}

func finalizeSearchMetadata(metadata *SearchMetadata) {
	if metadata.m != nil {
		C.destroy_search_metadata((*C._SearchMetadata)(metadata.m))
	}
	metadata.m = nil
}

func makeSearchMetadata(m *C._SearchMetadata) *SearchMetadata {
	metadata := new(SearchMetadata)
	runtime.SetFinalizer(metadata, finalizeSearchMetadata)
	metadata.m = (*C._QueryMetadata)(m)
	return metadata
}

// NewSearchMetadata creates a new SearchMetadata with the given locale and
// form_factor
func NewSearchMetadata(cardinality int, locale, form_factor string) *SearchMetadata {
	return makeSearchMetadata(C.new_search_metadata(C.int(cardinality),
		unsafe.Pointer(&locale),
		unsafe.Pointer(&form_factor)))
}

// Cardinality returns the desired number of results for the search request.
func (metadata *SearchMetadata) Cardinality() int {
	return int(C.search_metadata_get_cardinality((*C._SearchMetadata)(metadata.m)))
}

type Location struct {
	Latitude           float64 `json:"latitude"`
	Longitude          float64 `json:"longitude"`
	Altitude           float64 `json:"altitude"`
	AreaCode           string  `json:"area_code"`
	City               string  `json:"city"`
	CountryCode        string  `json:"country_code"`
	CountryName        string  `json:"country_name"`
	HorizontalAccuracy float64 `json:"horizontal_accuracy"`
	VerticalAccuracy   float64 `json:"vertical_accuracy"`
	RegionCode         string  `json:"region_code"`
	RegionName         string  `json:"region_name"`
	ZipPostalCode      string  `json:"zip_postal_code"`
}

func (metadata *SearchMetadata) Location() *Location {
	var length C.int
	locData := C.search_metadata_get_location((*C._SearchMetadata)(metadata.m), &length)
	if locData == nil {
		return nil
	}
	defer C.free(locData)
	var location Location
	if err := json.Unmarshal(C.GoBytes(locData, length), &location); err != nil {
		panic(err)
	}
	return &location
}

// SetLocation sets the location
func (metadata *SearchMetadata) SetLocation(l *Location) error {
	location := locationMarshal{marshalFloat(l.Latitude),
		marshalFloat(l.Longitude),
		marshalFloat(l.Altitude),
		l.AreaCode,
		l.City,
		l.CountryCode,
		l.CountryName,
		marshalFloat(l.HorizontalAccuracy),
		marshalFloat(l.VerticalAccuracy),
		l.RegionCode,
		l.RegionName,
		l.ZipPostalCode}
	data, err := json.Marshal(location)
	if err != nil {
		return err
	}
	var errorString *C.char
	C.search_metadata_set_location((*C._SearchMetadata)(metadata.m), (*C.char)(unsafe.Pointer(&data[0])), C.int(len(data)), &errorString)
	return checkError(errorString)
}

func (metadata *SearchMetadata) SetAggregatedKeywords(keywords []string) error {
	var errorString *C.char
	C.search_metadata_set_aggregated_keywords((*C._SearchMetadata)(metadata.m), unsafe.Pointer(&keywords[0]), C.int(len(keywords)), &errorString)
	return checkError(errorString)
}

func (metadata *SearchMetadata) AggregatedKeywords() []string {
	var length C.int
	keywordData := C.search_metadata_get_aggregated_keywords((*C._SearchMetadata)(metadata.m), &length)
	var keywords []string
	if err := json.Unmarshal(C.GoBytes(keywordData, length), &keywords); err != nil {
		panic(err)
	}
	return keywords
}

func (metadata *SearchMetadata) IsAggregated() bool {
	if C.search_metadata_is_aggregated((*C._SearchMetadata)(metadata.m)) == 0 {
		return false
	}
	return true
}

// ActionMetadata holds additional metadata about the preview request
// or result activation.
type ActionMetadata struct {
	queryMetadata
}

func finalizeActionMetadata(metadata *ActionMetadata) {
	if metadata.m != nil {
		C.destroy_action_metadata((*C._ActionMetadata)(metadata.m))
	}
	metadata.m = nil
}

// NewActionMetadata creates a new ActionMetadata with the given locale and
// form_factor
func NewActionMetadata(locale, form_factor string) *ActionMetadata {
	return makeActionMetadata(C.new_action_metadata(unsafe.Pointer(&locale),
		unsafe.Pointer(&form_factor)))
}

func makeActionMetadata(m *C._ActionMetadata) *ActionMetadata {
	metadata := new(ActionMetadata)
	runtime.SetFinalizer(metadata, finalizeActionMetadata)
	metadata.m = (*C._QueryMetadata)(m)
	return metadata
}

// ScopeData decodes the stored scope data into the given variable.
//
// Scope data is either set by the shell when calling a preview
// action, or set by the scope through an ActivationResponse object.
func (metadata *ActionMetadata) ScopeData(v interface{}) error {
	var dataLength C.int
	scopeData := C.action_metadata_get_scope_data((*C._ActionMetadata)(metadata.m), &dataLength)
	defer C.free(scopeData)
	return json.Unmarshal(C.GoBytes(scopeData, dataLength), v)
}

// SetScopeData attaches arbitrary data to this ActionMetadata.
func (metadata *ActionMetadata) SetScopeData(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var errorString *C.char
	C.action_metadata_set_scope_data((*C._ActionMetadata)(metadata.m), (*C.char)(unsafe.Pointer(&data[0])), C.int(len(data)), &errorString)
	return checkError(errorString)
}

// SetHint sets a hint.
func (metadata *ActionMetadata) SetHint(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var errorString *C.char
	C.action_metadata_set_hint((*C._ActionMetadata)(metadata.m), unsafe.Pointer(&key), (*C.char)(unsafe.Pointer(&data[0])), C.int(len(data)), &errorString)
	return checkError(errorString)
}

// Hint returns a hint.
// Returns error if the hint does not exist or if we got an error unmarshaling
func (metadata *ActionMetadata) Hint(key string, value interface{}) error {
	var dataLength C.int
	var errorString *C.char
	scopeData := C.action_metadata_get_hint((*C._ActionMetadata)(metadata.m), unsafe.Pointer(&key), &dataLength, &errorString)
	if dataLength > 0 && errorString == nil {
		defer C.free(scopeData)
		return json.Unmarshal(C.GoBytes(scopeData, dataLength), value)
	} else {
		return checkError(errorString)
	}
}

// Hints gets all hints.
func (metadata *ActionMetadata) Hints(value interface{}) error {
	var length C.int
	data := C.action_metadata_get_hints((*C._ActionMetadata)(metadata.m), &length)
	if data == nil {
		return nil
	}
	defer C.free(data)
	return json.Unmarshal(C.GoBytes(data, length), value)
}

// we use this type to reimplement the marshaller interface in order to make values
// like 1.0 not being converted as 1 (integer).
type marshalFloat float64

func (n marshalFloat) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%f", n)), nil
}

// the following structure is only used to control how the float64 types are
// marshaled. It is not exported.
type locationMarshal struct {
	Latitude           marshalFloat `json:"latitude"`
	Longitude          marshalFloat `json:"longitude"`
	Altitude           marshalFloat `json:"altitude"`
	AreaCode           string       `json:"area_code"`
	City               string       `json:"city"`
	CountryCode        string       `json:"country_code"`
	CountryName        string       `json:"country_name"`
	HorizontalAccuracy marshalFloat `json:"horizontal_accuracy"`
	VerticalAccuracy   marshalFloat `json:"vertical_accuracy"`
	RegionCode         string       `json:"region_code"`
	RegionName         string       `json:"region_name"`
	ZipPostalCode      string       `json:"zip_postal_code"`
}

type ProxyScopeMetadata struct {
	Identity string `json:"identity"`
	EndPoint string `json:"endpoint"`
}

// ScopeMetadata holds scope attributes such as name, description, icon etc.
//
// The information stored by ScopeMetadata comes from the .ini file for the given scope (for local scopes)
// or is fetched from the remote server (for scopes running on Smart Scopes Server).
// Use ListRegistryScopes from ScopeBase to get the metadata for all scopes.
type ScopeMetadata struct {
	m                    *C._ScopeMetadata
	Art                  string                 `json:"art"`
	Author               string                 `json:"author"`
	Description          string                 `json:"description"`
	DisplayName          string                 `json:"display_name"`
	Icon                 string                 `json:"icon"`
	Invisible            bool                   `json:"invisible"`
	IsAggregator         bool                   `json:"is_aggregator"`
	LocationDataNeeded   bool                   `json:"location_data_needed"`
	ScopeDir             string                 `json:"scope_dir"`
	ScopeId              string                 `json:"scope_id"`
	Version              int                    `json:"version"`
	Proxy                ProxyScopeMetadata     `json:"proxy"`
	AppearanceAttributes map[string]interface{} `json:"appearance_attributes"`
	SettingsDefinitions  []interface{}          `json:"settings_definitions"`
	Keywords             []string               `json:"keywords"`
}

func finalizeScopeMetadata(metadata *ScopeMetadata) {
	C.destroy_scope_metadata_ptr(metadata.m)
}

func makeScopeMetadata(m *C._ScopeMetadata, json_data string) *ScopeMetadata {
	metadata := new(ScopeMetadata)
	if err := json.Unmarshal([]byte(json_data), &metadata); err != nil {
		panic(err)
	}
	metadata.m = m
	runtime.SetFinalizer(metadata, finalizeScopeMetadata)
	return metadata
}
