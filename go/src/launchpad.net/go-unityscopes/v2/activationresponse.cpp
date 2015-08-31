#include <stdexcept>
#include <cstring>

#include <unity/scopes/ActivationResponse.h>
#include <unity/scopes/CannedQuery.h>

extern "C" {
#include "_cgo_export.h"
}

using namespace unity::scopes;

void activation_response_init_status(_ActivationResponse *response, int status) {
    *reinterpret_cast<ActivationResponse*>(response) =
        ActivationResponse(static_cast<ActivationResponse::Status>(status));
}

void activation_response_init_query(_ActivationResponse *response, _CannedQuery *query) {
    *reinterpret_cast<ActivationResponse*>(response) =
        ActivationResponse(*reinterpret_cast<CannedQuery*>(query));
}

void activation_response_set_scope_data(_ActivationResponse *response, char *json_data, int json_data_length, char **error) {
    try {
        Variant v = Variant::deserialize_json(std::string(json_data, json_data_length));
        reinterpret_cast<ActivationResponse*>(response)->set_scope_data(v);
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}
