#include "helpers.h"

#include <cstdlib>
#include <cstring>

extern "C" {
#include "_cgo_export.h"
}

namespace gounityscopes {
namespace internal {

std::string from_gostring(void *str) {
    GoString *s = static_cast<GoString*>(str);
    return std::string(s->p, s->n);
}

void *as_bytes(const std::string &str, int *length) {
    *length = str.size();
    void *data = malloc(str.size());
    if (data == nullptr) {
        return nullptr;
    }
    memcpy(data, str.data(), str.size());
    return data;
}

}
}
