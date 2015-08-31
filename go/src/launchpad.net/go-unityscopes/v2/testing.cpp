#include <stdexcept>
#include <cstring>

#include <unity/scopes/testing/Result.h>

extern "C" {
#include "_cgo_export.h"
}
#include "helpers.h"

using namespace unity::scopes;
using namespace gounityscopes::internal;

_Result *new_testing_result() {
    return reinterpret_cast<_Result*>(static_cast<Result*>(new testing::Result));
}
