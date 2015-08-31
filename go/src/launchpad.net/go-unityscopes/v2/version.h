#ifndef UNITYSCOPE_GO_VERSION_H
#define UNITYSCOPE_GO_VERSION_H

#include <unity/scopes/Version.h>

// check that we have a compatible version of lib-unityscopes installed
static_assert(UNITY_SCOPES_VERSION_MAJOR > 0 ||
             (UNITY_SCOPES_VERSION_MAJOR == 0 && (UNITY_SCOPES_VERSION_MINOR > 6 || (UNITY_SCOPES_VERSION_MINOR == 6 && UNITY_SCOPES_VERSION_MICRO >= 15))),
              "Version of Unity scopes API mismatch. Minimum required version is 0.6.15.");

#endif
