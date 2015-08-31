#include <stdexcept>
#include <cstring>

#include <unity/scopes/CannedQuery.h>

extern "C" {
#include "_cgo_export.h"
}
#include "helpers.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

void destroy_canned_query(_CannedQuery *query) {
    delete reinterpret_cast<CannedQuery*>(query);
}

_CannedQuery *new_canned_query(void *scope_id, void *query_str, void *department_id) {
    return reinterpret_cast<_CannedQuery*>(
        new CannedQuery(from_gostring(scope_id),
                        from_gostring(query_str),
                        from_gostring(department_id)));
}

char *canned_query_get_scope_id(_CannedQuery *query) {
    return strdup(reinterpret_cast<CannedQuery*>(query)->scope_id().c_str());
}

char *canned_query_get_department_id(_CannedQuery *query) {
    return strdup(reinterpret_cast<CannedQuery*>(query)->department_id().c_str());
}

void *canned_query_get_filter_state(_CannedQuery *query, int *length) {
    std::string json_data;
    try {
        Variant v(reinterpret_cast<CannedQuery*>(query)->filter_state().serialize());
        json_data = v.serialize_json();
    } catch (...) {
        return nullptr;
    }
    return as_bytes(json_data, length);
}

char *canned_query_get_query_string(_CannedQuery *query) {
    return strdup(reinterpret_cast<CannedQuery*>(query)->query_string().c_str());
}

void canned_query_set_department_id(_CannedQuery *query, void *department_id) {
    reinterpret_cast<CannedQuery*>(query)->set_department_id(from_gostring(department_id));
}

void canned_query_set_query_string(_CannedQuery *query, void *query_str) {
    reinterpret_cast<CannedQuery*>(query)->set_query_string(from_gostring(query_str));
}

char *canned_query_to_uri(_CannedQuery *query) {
    return strdup(reinterpret_cast<CannedQuery*>(query)->to_uri().c_str());
}
