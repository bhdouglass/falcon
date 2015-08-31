#include <stdexcept>
#include <cstring>

#include <unity/scopes/Department.h>
#include <unity/scopes/CannedQuery.h>

extern "C" {
#include "_cgo_export.h"
}
#include "helpers.h"
#include "smartptr_helper.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

/* Department objects */
void init_department_ptr(SharedPtrData dest, SharedPtrData src) {
    std::shared_ptr<Department> dept = get_ptr<Department>(src);
    init_ptr<Department>(dest, dept);
}

void new_department(void *dept_id, _CannedQuery *query, void *label, SharedPtrData dept, char **error) {
    try {
        Department::UPtr d;
        if(dept_id) {
            d = Department::create(from_gostring(dept_id),
                                    *reinterpret_cast<CannedQuery*>(query),
                                    from_gostring(label));
        } else {
            d = Department::create(*reinterpret_cast<CannedQuery*>(query),
                                   from_gostring(label));
        }
        init_ptr<Department>(dept, std::move(d));
    } catch (const std::exception &e) {
        *error = strdup(e.what());
    }
}

void destroy_department_ptr(SharedPtrData data) {
    destroy_ptr<Department>(data);
}

void department_add_subdepartment(SharedPtrData dept, SharedPtrData child) {
    get_ptr<Department>(dept)->add_subdepartment(get_ptr<Department>(child));
}

void department_set_alternate_label(SharedPtrData dept, void *label) {
    get_ptr<Department>(dept)->set_alternate_label(from_gostring(label));
}

char *department_get_alternate_label(SharedPtrData dept) {
    return strdup(get_ptr<Department>(dept)->alternate_label().c_str());
}

char *department_get_id(SharedPtrData dept) {
    return strdup(get_ptr<Department>(dept)->id().c_str());
}

char *department_get_label(SharedPtrData dept) {
    return strdup(get_ptr<Department>(dept)->label().c_str());
}

void department_set_has_subdepartments(SharedPtrData dept, int subdepartments) {
    get_ptr<Department>(dept)->set_has_subdepartments(subdepartments);
}

int department_has_subdepartments(SharedPtrData dept) {
    return static_cast<int>(get_ptr<Department>(dept)->has_subdepartments());
}

SharedPtrData * department_get_subdepartments(SharedPtrData dept, int *n_subdepts) {
    auto subdepts = get_ptr<Department>(dept)->subdepartments();
    *n_subdepts = subdepts.size();
    SharedPtrData* ret_data =
    reinterpret_cast<SharedPtrData*>(calloc(*n_subdepts, sizeof(SharedPtrData)));
    int i = 0;
    for(auto item: subdepts) {
        init_ptr<Department const>(ret_data[i++], item);
    }
    return ret_data;
}

void department_set_subdepartments(SharedPtrData dept, SharedPtrData **subdepartments, int nb_subdepartments) {
    DepartmentList api_depts;
    for(auto i = 0; i < nb_subdepartments; i++) {
        api_depts.push_back(get_ptr<Department>(*subdepartments[i]));
    }
    get_ptr<Department>(dept)->set_subdepartments(api_depts);
}

_CannedQuery * department_get_query(SharedPtrData dept) {
    return reinterpret_cast<_CannedQuery*>(new CannedQuery(get_ptr<Department>(dept)->query()));
}
