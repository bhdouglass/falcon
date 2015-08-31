#ifndef UNITYSCOPE_SMARTPTR_HELPER_H
#define UNITYSCOPE_SMARTPTR_HELPER_H

#include <memory>

#include "shim.h"

namespace gounityscopes {
namespace internal {

template<typename T> inline void init_const_ptr(SharedPtrData data,  std::shared_ptr<T const> v) {
    typedef std::shared_ptr<T const> Ptr;
    static_assert(sizeof(SharedPtrData) >= sizeof(Ptr),
                  "std::shared_ptr is larger than expected");
    Ptr *ptr = new(reinterpret_cast<void*>(data)) Ptr();
    *ptr = v;
}

template<typename T> inline void init_ptr(SharedPtrData data, std::shared_ptr<T> v) {
    typedef std::shared_ptr<T> Ptr;
    static_assert(sizeof(SharedPtrData) >= sizeof(Ptr),
                  "std::shared_ptr is larger than expected");
    Ptr *ptr = new(reinterpret_cast<void*>(data)) Ptr();
    *ptr = v;
}

template<typename T> inline std::shared_ptr<T> get_ptr(SharedPtrData data) {
    typedef std::shared_ptr<T> Ptr;
    Ptr *ptr = reinterpret_cast<Ptr*>(data);
    return *ptr;
}

template<typename T> inline void destroy_ptr(SharedPtrData data) {
    typedef std::shared_ptr<T> Ptr;
    if (!(data[0] == 0 && data[1] == 0)) {
        Ptr *ptr = reinterpret_cast<Ptr*>(data);
        ptr->~Ptr();
    }
    data[0] = data[1] = 0;
}

template<typename T> inline void copy_ptr(SharedPtrData dest_data, SharedPtrData src_data) {
    typedef std::shared_ptr<T> Ptr;
    Ptr *dest = reinterpret_cast<Ptr*>(dest_data);
    Ptr *src = reinterpret_cast<Ptr*>(src_data);
    *dest = *src;
}

}
}

#endif
