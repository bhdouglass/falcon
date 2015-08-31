#include <stdexcept>
#include <cstring>

#include <unity/scopes/Category.h>
#include <unity/scopes/CategorisedResult.h>
#include <unity/scopes/Result.h>

extern "C" {
#include "_cgo_export.h"
}
#include "helpers.h"
#include "smartptr_helper.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

_Result *new_categorised_result(SharedPtrData category) {
    auto cat = get_ptr<Category>(category);
    return reinterpret_cast<_CategorisedResult*>(static_cast<Result*>(new CategorisedResult(cat)));
}

void destroy_result(_Result *res) {
    delete reinterpret_cast<Result*>(res);
}

void *result_get_attr(_Result *res, void *attr, int *length, char **error) {
    std::string json_data;
    try {
        Variant v = reinterpret_cast<Result*>(res)->value(from_gostring(attr));
        json_data = v.serialize_json();
    } catch (const std::exception &e) {
        *error = strdup(e.what());
        return nullptr;
    }
    return as_bytes(json_data, length);
}

void result_set_attr(_Result *res, void *attr, void *json_value, char **error) {
    try {
        Variant v = Variant::deserialize_json(from_gostring(json_value));
        (*reinterpret_cast<Result*>(res))[from_gostring(attr)] = v;
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}

void result_set_intercept_activation(_Result *res) {
    reinterpret_cast<Result*>(res)->set_intercept_activation();
}
