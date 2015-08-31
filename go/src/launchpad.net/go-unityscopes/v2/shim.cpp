#include <cstring>

#include <unity/scopes/Category.h>
#include <unity/scopes/Runtime.h>

extern "C" {
#include "_cgo_export.h"
}
#include "version.h"
#include "helpers.h"
#include "smartptr_helper.h"
#include "scope.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

void run_scope(void *scope_name, void *runtime_config, void *scope_config,
               void *pointer_to_iface, char **error) {
    try {
        auto runtime = Runtime::create_scope_runtime(
            from_gostring(scope_name), from_gostring(runtime_config));
        ScopeAdapter scope(*reinterpret_cast<GoInterface*>(pointer_to_iface));
        runtime->run_scope(&scope, from_gostring(scope_config));
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}

char *scope_base_scope_directory(_ScopeBase *scope) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    return strdup(s->scope_directory().c_str());
}

char *scope_base_cache_directory(_ScopeBase *scope) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    return strdup(s->cache_directory().c_str());
}

char *scope_base_tmp_directory(_ScopeBase *scope) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    return strdup(s->tmp_directory().c_str());
}

void *scope_base_settings(_ScopeBase *scope, int *length) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    Variant settings(s->settings());
    return as_bytes(settings.serialize_json(), length);
}

void destroy_category_ptr(SharedPtrData data) {
    destroy_ptr<const Category>(data);
}

_ScopeMetadata** list_registry_scopes_metadata(_ScopeBase *scope, int *n_scopes) {
    ScopeBase *s = reinterpret_cast<ScopeBase*>(scope);
    auto registry = s->registry();
    auto scopes = registry->list();

    *n_scopes = scopes.size();

    _ScopeMetadata** ret_data = reinterpret_cast<_ScopeMetadata**>(calloc(*n_scopes, sizeof(_ScopeMetadata*)));
    int i = 0;
    for( auto item: scopes ) {
        ret_data[i++] = reinterpret_cast<_ScopeMetadata*>(new ScopeMetadata(item.second));
    }

    return ret_data;
}
