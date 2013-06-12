#ifndef _minunit_h
#define _minunit_h

#include <stdbool.h>

//==============================================================================
//
// Minunit
//
//==============================================================================

#define mu_fail(MSG, ...) do {\
    fprintf(stderr, "%s:%d: " MSG "\n", __FILE__, __LINE__, ##__VA_ARGS__);\
    return 1;\
} while(0)

#define mu_assert_with_msg(TEST, MSG, ...) do {\
    if (!(TEST)) {\
        fprintf(stderr, "%s:%d: %s " MSG "\n", __FILE__, __LINE__, #TEST, ##__VA_ARGS__);\
        return 1;\
    }\
} while (0)

#define mu_assert(TEST, MSG, ...) do {\
    if (!(TEST)) {\
        fprintf(stderr, "%s:%d: %s " MSG "\n", __FILE__, __LINE__, #TEST, ##__VA_ARGS__);\
        return 1;\
    }\
} while (0)

#define mu_run_test(TEST) do {\
    fprintf(stderr, "%s\n", #TEST);\
    int rc = TEST();\
    tests_run++; \
    if (rc) {\
        fprintf(stderr, "\n  Test Failure: %s()\n", #TEST);\
        return rc;\
    }\
} while (0)

#define RUN_TESTS() int main() {\
   fprintf(stderr, "== %s ==\n", __FILE__);\
   int rc = all_tests();\
   fprintf(stderr, "\n");\
   return rc;\
}

int tests_run;


//--------------------------------------
// Typed asserts
//--------------------------------------

#define mu_assert_bool(ACTUAL) mu_assert_with_msg(ACTUAL, "Expected value to be true but was %d", ACTUAL)
#define mu_assert_int_equals(ACTUAL, EXPECTED) mu_assert_with_msg(ACTUAL == EXPECTED, "Expected: %d; Received: %d", EXPECTED, ACTUAL)
#define mu_assert_long_equals(ACTUAL, EXPECTED) mu_assert_with_msg(ACTUAL == EXPECTED, "Expected: %ld; Received: %ld", EXPECTED, ACTUAL)
#define mu_assert_int64_equals(ACTUAL, EXPECTED) mu_assert_with_msg(ACTUAL == EXPECTED, "Expected: %lld; Received: %lld", (long long int)EXPECTED, (long long int)ACTUAL)

#endif