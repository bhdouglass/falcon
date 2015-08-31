#include <stdexcept>
#include <cmath>
#include <cstring>

#include <unity/scopes/ActionMetadata.h>
#include <unity/scopes/SearchMetadata.h>
#include <unity/scopes/ScopeMetadata.h>
#include <unity/scopes/ScopeExceptions.h>

extern "C" {
#include "_cgo_export.h"
}
#include "helpers.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

/* SearchMetadata objects */
_SearchMetadata *new_search_metadata(int cardinality, void *locale, void *form_factor) {
    return reinterpret_cast<_SearchMetadata*>(new SearchMetadata(cardinality,
                                                                from_gostring(locale),
                                                                from_gostring(form_factor)));
}

void destroy_search_metadata(_SearchMetadata *metadata) {
    delete reinterpret_cast<SearchMetadata*>(metadata);
}

char *query_metadata_get_locale(_QueryMetadata *metadata) {
    auto m = reinterpret_cast<QueryMetadata*>(metadata);
    try {
        return strdup(m->locale().c_str());
    } catch (const NotFoundException &) {
        return nullptr;
    }
}

char *query_metadata_get_form_factor(_QueryMetadata *metadata) {
    auto m = reinterpret_cast<QueryMetadata*>(metadata);
    try {
        return strdup(m->form_factor().c_str());
    } catch (const NotFoundException &) {
        return nullptr;
    }
}

void query_metadata_set_internet_connectivity(_QueryMetadata *metadata, int status) {
    reinterpret_cast<QueryMetadata*>(metadata)->set_internet_connectivity(static_cast<QueryMetadata::ConnectivityStatus>(status));
}

int query_metadata_get_internet_connectivity(_QueryMetadata *metadata) {
    return static_cast<int>(reinterpret_cast<QueryMetadata*>(metadata)->internet_connectivity());
}

int search_metadata_get_cardinality(_SearchMetadata *metadata) {
    return reinterpret_cast<SearchMetadata*>(metadata)->cardinality();
}

void *search_metadata_get_location(_SearchMetadata *metadata, int *length) {
    auto m = reinterpret_cast<SearchMetadata*>(metadata);
    VariantMap location;
    try {
        location = m->location().serialize();
    } catch (const NotFoundException &) {
        return nullptr;
    }
    // libjsoncpp generates invalid JSON for NaN or Inf values, so
    // filter them out here.
    for (auto &pair : location) {
        if (pair.second.which() == Variant::Double) {
            double value = pair.second.get_double();
            if (!std::isfinite(value)) {
                pair.second = Variant();
            }
        }
    }
    return as_bytes(Variant(location).serialize_json(), length);
}

void search_metadata_set_location(_SearchMetadata *metadata, char *json_data, int json_data_length, char **error) {

    try {
        Variant value = Variant::deserialize_json(std::string(json_data, json_data_length));
        Location location(value.get_dict());
        reinterpret_cast<SearchMetadata*>(metadata)->set_location(location);
    } catch (const std::exception & e) {
        *error = strdup(e.what());
    }
}

void search_metadata_set_aggregated_keywords(_SearchMetadata *metadata, void *gostring_array, int count, char **error) {
    try {
        GoString *keyword_data = static_cast<GoString*>(gostring_array);
        std::set<std::string> keywords;
        for (int i = 0; i < count; i++) {
            keywords.emplace(std::string(keyword_data[i].p, keyword_data[i].n));
        }
        reinterpret_cast<SearchMetadata*>(metadata)->set_aggregated_keywords(keywords);
    } catch (const std::exception & e) {
        *error = strdup(e.what());
    }
}

void *search_metadata_get_aggregated_keywords(_SearchMetadata *metadata, int *length) {
    std::set<std::string> keywords = reinterpret_cast<SearchMetadata*>(metadata)->aggregated_keywords();
    // Marshal via JSON for now.  This is probably faster than calling
    // C.free() on each of a list of strings.
    VariantArray array(keywords.begin(), keywords.end());
    return as_bytes(Variant(array).serialize_json(), length);
}

int search_metadata_is_aggregated(_SearchMetadata *metadata) {
    return reinterpret_cast<SearchMetadata*>(metadata)->is_aggregated();
}


/* ActionMetadata objects */
_ActionMetadata *new_action_metadata(void *locale, void *form_factor) {
    return reinterpret_cast<_ActionMetadata*>(new ActionMetadata(from_gostring(locale),
                                                                 from_gostring(form_factor)));
}

void destroy_action_metadata(_ActionMetadata *metadata) {
    delete reinterpret_cast<ActionMetadata*>(metadata);
}

void *action_metadata_get_scope_data(_ActionMetadata *metadata, int *data_length) {
    const std::string data = reinterpret_cast<ActionMetadata*>(metadata)->scope_data().serialize_json();
    return as_bytes(data, data_length);
}

void action_metadata_set_scope_data(_ActionMetadata *metadata, char *json_data, int json_data_length, char **error) {
    try {
        Variant value = Variant::deserialize_json(std::string(json_data, json_data_length));
        reinterpret_cast<ActionMetadata*>(metadata)->set_scope_data(value);
    } catch (const std::exception & e) {
        *error = strdup(e.what());
    }
}

void action_metadata_set_hint(_ActionMetadata *metadata, void *key, char *json_data, int json_data_length, char **error) {
    try {
        Variant value = Variant::deserialize_json(std::string(json_data, json_data_length));
        reinterpret_cast<ActionMetadata*>(metadata)->set_hint(from_gostring(key), value);
    } catch (const std::exception & e) {
        *error = strdup(e.what());
    }
}

void *action_metadata_get_hint(_ActionMetadata *metadata, void *key, int *data_length, char **error) {
    try {
        ActionMetadata const*api_metadata = reinterpret_cast<ActionMetadata const*>(metadata);
        Variant value = (*api_metadata)[from_gostring(key)];
        const std::string data = value.serialize_json();
        return as_bytes(data, data_length);
    } catch (const std::exception & e) {
        *data_length = 0;
        *error = strdup(e.what());
        return 0;
    }
}

void *action_metadata_get_hints(_ActionMetadata *metadata, int *length) {
    VariantMap hints = reinterpret_cast<ActionMetadata const*>(metadata)->hints();
    // libjsoncpp generates invalid JSON for NaN or Inf values, so
    // filter them out here.
    for (auto &pair : hints) {
        if (pair.second.which() == Variant::Double) {
            double value = pair.second.get_double();
            if (!std::isfinite(value)) {
                pair.second = Variant();
            }
        }
    }
    return as_bytes(Variant(hints).serialize_json(), length);
}

char *get_scope_metadata_serialized(_ScopeMetadata *metadata) {
    ScopeMetadata const*api_metadata = reinterpret_cast<ScopeMetadata const*>(metadata);
    return strdup(Variant(api_metadata->serialize()).serialize_json().c_str());
}

void destroy_scope_metadata_ptr(_ScopeMetadata *metadata) {
    delete reinterpret_cast<ScopeMetadata*>(metadata);
}
