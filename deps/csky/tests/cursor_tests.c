#include <stdio.h>
#include <stdlib.h>
#include <stddef.h>
#include <math.h>

#include <sky/cursor.h>
#include <sky/sky_string.h>
#include <sky/timestamp.h>
#include <sky/mem.h>

#include "minunit.h"

//==============================================================================
//
// Fixtures
//
//==============================================================================

char INT_DATA[] = "\xD1\x03\xE8";

char DOUBLE_DATA[] = "\xCB\x40\x59\x0C\xCC\xCC\xCC\xCC\xCD";

char BOOLEAN_FALSE_DATA[] = "\xC2";

char BOOLEAN_TRUE_DATA[] = "\xC3";

char STRING_DATA[] = "\xa3\x66\x6f\x6f";


int DATA0_LENGTH = 129;
char *DATA0 = "\xA0"
  // 1970-01-01T00:00:00Z, {1:"john doe", 2:1000, 3:100.2, 4:true}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x00\x00\x00" "\x84" "\x01\xA8""john doe" "\x02\xD1\x03\xE8" "\x03\xCB\x40\x59\x0C\xCC\xCC\xCC\xCC\xCD" "\x04\xC3"
  // 1970-01-01T00:00:01Z, {-1:"A1", -2:"super", -3:21, -4:100, -5:true}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x10\x00\x00" "\x85" "\xFF\xA2""A1" "\xFE\xA5""super" "\xFD\x15" "\xFC\xCB\x40\x59\x00\x00\x00\x00\x00\x00" "\xFB\xC3"
  // 1970-01-01T00:00:02Z, {-1:"A2"}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x20\x00\x00" "\x81" "\xFF\xA2""A2"
  // 1970-01-01T00:00:03Z, {1:"frank sinatra", 2:20, 3:-100, 4:false}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x30\x00\x00" "\x84" "\x01\xAD""frank sinatra" "\x02\x14" "\x03\xCB\xC0\x59\x00\x00\x00\x00\x00\x00" "\x04\xC2"
;

int DATA1_LENGTH = 112;
char *DATA1 = "\xA0"
  // 1970-01-01T00:00:00Z, {-1:"A1", 1:1000}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x00\x00\x00" "\x82" "\xFF\xA2""A1" "\x01\xD1\x03\xE8"
  // 1970-01-01T00:00:01Z, {-1:"A2", -2:100}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x10\x00\x00" "\x82" "\xFF\xA2""A2" "\xFE\x64"
  // 1970-01-01T00:00:10Z, {-1:"A3", -2:200}
  "\x92" "\xD3\x00\x00\x00\x00\x00\xA0\x00\x00" "\x82" "\xFF\xA2""A3" "\xFE\xD1\x00\xC8"
  // 1970-01-01T00:00:20Z, {-1:"A1", -2:300}
  "\x92" "\xD3\x00\x00\x00\x00\x01\x40\x00\x00" "\x82" "\xFF\xA2""A1" "\xFE\xD1\x01\x2C"
  // 1970-01-01T00:01:00Z, {-1:"A1", 1:2000}
  "\x92" "\xD3\x00\x00\x00\x00\x03\xC0\x00\x00" "\x82" "\xFF\xA2""A1" "\x01\xD1\x07\xD0"
  // 1970-01-01T00:01:00Z, {-1:"A1", 1:2000}
  "\x92" "\xD3\x00\x00\x00\x00\x03\xF0\x00\x00" "\x82" "\xFF\xA2""A2" "\xFE\xD1\x01\x90"
;

int DATA3_LENGTH = 27;
char *DATA3 = "\xA0"
  // 1970-01-01T00:00:00Z, {1:2}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x00\x00\x00" "\x81" "\x01\x02"
  // 1970-01-01T00:00:01Z, {1:3}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x10\x00\x00" "\x81" "\x01\x03"
;

int DATA4_LENGTH = 14;
char *DATA4 = "\xA0"
  // 1970-01-01T00:00:00Z, {1:2}
  "\x92" "\xD3\x00\x00\x00\x00\x00\x00\x00\x00" "\x81" "\x01\x04"
;

int DATA5_LENGTH = 14;
char *DATA5 = "\xA0"
  // 1970-01-01T00:00:10Z, {1:20}
  "\x92" "\xD3\x00\x00\x00\x00\x00\xA0\x00\x00" "\x81" "\x01\x14"
;


//==============================================================================
//
// Declarations
//
//==============================================================================

#define ASSERT_OBJ_STATE(OBJ, TS, TIMESTAMP, ACTION, OSTRING, OINT, ODOUBLE, OBOOLEAN, ASTRING, AINT, ADOUBLE, ABOOLEAN) do {\
    mu_assert_int64_equals(((test_t*)OBJ)->ts, TS); \
    mu_assert_int_equals(((test_t*)OBJ)->timestamp, TIMESTAMP); \
    mu_assert_int_equals(((test_t*)OBJ)->action.length, (int)strlen(ACTION)); \
    mu_assert_bool(memcmp(((test_t*)OBJ)->action.data, ACTION, strlen(ACTION)) == 0); \
    mu_assert_int_equals(((test_t*)OBJ)->object_string.length, (int)strlen(OSTRING)); \
    mu_assert_bool(memcmp(((test_t*)OBJ)->object_string.data, OSTRING, strlen(OSTRING)) == 0); \
    mu_assert_int64_equals(((test_t*)OBJ)->object_int, OINT); \
    mu_assert_bool(fabs(((test_t*)OBJ)->object_double-ODOUBLE) < 0.1); \
    mu_assert_bool(((test_t*)OBJ)->object_boolean == OBOOLEAN); \
    mu_assert_int_equals(((test_t*)OBJ)->action_string.length, (int)strlen(ASTRING)); \
    mu_assert_bool(memcmp(((test_t*)OBJ)->action_string.data, ASTRING, strlen(ASTRING)) == 0); \
    mu_assert_int64_equals(((test_t*)OBJ)->action_int, AINT); \
    mu_assert_bool(fabs(((test_t*)OBJ)->action_double-ADOUBLE) < 0.1); \
    mu_assert_bool(((test_t*)OBJ)->action_boolean == ABOOLEAN); \
} while(0)

#define ASSERT_OBJ_STATE2(OBJ, TIMESTAMP, ACTION, OINT, AINT) do {\
    mu_assert_int64_equals(((test_t*)OBJ)->timestamp, TIMESTAMP); \
    mu_assert_int_equals(((test_t*)OBJ)->action.length, (int)strlen(ACTION)); \
    mu_assert_bool(memcmp(((test_t*)OBJ)->action.data, ACTION, strlen(ACTION)) == 0); \
    mu_assert_int64_equals(((test_t*)OBJ)->object_int, OINT); \
    mu_assert_int64_equals(((test_t*)OBJ)->action_int, AINT); \
} while(0)

typedef struct {
    sky_string action;
    sky_string action_string;
    int32_t    action_int;
    double     action_double;
    bool       action_boolean;
    sky_string object_string;
    int32_t    object_int;
    double     object_double;
    bool       object_boolean;
    uint32_t timestamp;
    int64_t ts;
} test_t;

typedef struct {
    int64_t dummy;
    int32_t int_value;
    double double_value;
    bool boolean_value;
    sky_string string_value;
    uint32_t timestamp;
    int64_t ts;
} test2_t;

//==============================================================================
//
// Test Cases
//
//==============================================================================

//--------------------------------------
// Set Data
//--------------------------------------

int test_sky_cursor_set_data() {
    // Setup data object & cursor.
    sky_cursor *cursor = sky_cursor_new(-4, 4);
    sky_cursor_set_timestamp_offset(cursor, offsetof(test_t, timestamp));
    sky_cursor_set_ts_offset(cursor, offsetof(test_t, ts));
    sky_cursor_set_property(cursor, -5, offsetof(test_t, action_boolean), sizeof(bool), "boolean");
    sky_cursor_set_property(cursor, -4, offsetof(test_t, action_double), sizeof(double), "float");
    sky_cursor_set_property(cursor, -3, offsetof(test_t, action_int), sizeof(int32_t), "integer");
    sky_cursor_set_property(cursor, -2, offsetof(test_t, action_string), sizeof(sky_string), "string");
    sky_cursor_set_property(cursor, -1, offsetof(test_t, action), sizeof(sky_string), "string");
    sky_cursor_set_property(cursor, 1, offsetof(test_t, object_string), sizeof(sky_string), "string");
    sky_cursor_set_property(cursor, 2, offsetof(test_t, object_int), sizeof(int32_t), "integer");
    sky_cursor_set_property(cursor, 3, offsetof(test_t, object_double), sizeof(double), "float");
    sky_cursor_set_property(cursor, 4, offsetof(test_t, object_boolean), sizeof(bool), "boolean");
    sky_cursor_set_data_sz(cursor, sizeof(test_t));

    sky_cursor_set_ptr(cursor, DATA0, DATA0_LENGTH);
    ASSERT_OBJ_STATE(cursor->data, 0LL, 0, "", "", 0LL, 0, false, "", 0LL, 0, false);

    // Event 1 (State-Only)
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    ASSERT_OBJ_STATE(cursor->data, 0LL, 0, "", "john doe", 1000LL, 100.2, true, "", 0LL, 0, false);
    
    // Event 2 (Action + Action Data)
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    ASSERT_OBJ_STATE(cursor->data, sky_timestamp_shift(1000000LL), 1, "A1", "john doe", 1000LL, 100.2, true, "super", 21LL, 100, true);
    
    // Event 3 (Action-Only)
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    ASSERT_OBJ_STATE(cursor->data, sky_timestamp_shift(2000000LL), 2, "A2", "john doe", 1000LL, 100.2, true, "", 0LL, 0, false);

    // Event 4 (Data-Only)
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    ASSERT_OBJ_STATE(cursor->data, sky_timestamp_shift(3000000LL), 3, "", "frank sinatra", 20LL, -100, false, "", 0LL, 0, false);

    // EOF
    mu_assert_bool(!sky_lua_cursor_next_event(cursor));

    sky_cursor_free(cursor);
    return 0;
}


//--------------------------------------
// Sessionize
//--------------------------------------

int test_sky_cursor_sessionize() {
    // Setup data object.
    sky_cursor *cursor = sky_cursor_new(-2, 1);
    sky_cursor_set_timestamp_offset(cursor, offsetof(test_t, timestamp));
    sky_cursor_set_ts_offset(cursor, offsetof(test_t, ts));
    sky_cursor_set_property(cursor, -2, offsetof(test_t, action_int), sizeof(int32_t), "integer");
    sky_cursor_set_property(cursor, -1, offsetof(test_t, action), sizeof(sky_string), "string");
    sky_cursor_set_property(cursor, 1, offsetof(test_t, object_int), sizeof(int32_t), "integer");
    sky_cursor_set_data_sz(cursor, sizeof(test_t));

    // Initialize data and set a 10 second idle time.
    sky_cursor_set_ptr(cursor, DATA1, DATA1_LENGTH);
    sky_cursor_set_session_idle(cursor, 10);
    mu_assert_int_equals(cursor->session_event_index, -1);
    ASSERT_OBJ_STATE2(cursor->data, 0, "", 0LL, 0LL);
    
    // Pre-session
    mu_assert_bool(sky_lua_cursor_next_event(cursor) == false);
    mu_assert_int_equals(cursor->session_event_index, -1);
    ASSERT_OBJ_STATE2(cursor->data, 0, "", 0LL, 0LL);

    // Session 1
    mu_assert_bool(sky_lua_cursor_next_session(cursor));
    mu_assert_int_equals(cursor->session_event_index, -1);
    ASSERT_OBJ_STATE2(cursor->data, 0, "", 0LL, 0LL);

    // Session 1, Event 1
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 0);
    ASSERT_OBJ_STATE2(cursor->data, 0, "A1", 1000LL, 0LL);
    
    // Session 1, Event 2
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 1);
    ASSERT_OBJ_STATE2(cursor->data, 1, "A2", 1000LL, 100LL);
    
    // Session 1, Event 3
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 2);
    ASSERT_OBJ_STATE2(cursor->data, 10, "A3", 1000LL, 200LL);
    
    // Prevent next session!
    mu_assert_bool(sky_lua_cursor_next_event(cursor) == false);
    mu_assert_int_equals(cursor->session_event_index, 2);
    ASSERT_OBJ_STATE2(cursor->data, 10, "A3", 1000LL, 200LL);
    

    // Session 2 (Single Event)
    mu_assert_bool(sky_lua_cursor_next_session(cursor));
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 0);
    ASSERT_OBJ_STATE2(cursor->data, 20, "A1", 1000LL, 300LL);
    mu_assert_bool(sky_lua_cursor_next_event(cursor) == false);


    // Session 3 (with same data)
    mu_assert_bool(sky_lua_cursor_next_session(cursor));
    mu_assert_int_equals(cursor->session_event_index, -1);
    ASSERT_OBJ_STATE2(cursor->data, 20, "A1", 1000LL, 300LL);

    // Session 3, Event 1
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 0);
    ASSERT_OBJ_STATE2(cursor->data, 60, "A1", 2000LL, 0LL);

    // Session 3, Event 2
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 1);
    ASSERT_OBJ_STATE2(cursor->data, 63, "A2", 2000LL, 400LL);

    // Prevent next session!
    mu_assert_bool(sky_lua_cursor_next_event(cursor) == false);
    mu_assert_bool(sky_lua_cursor_next_session(cursor) == false);

    // EOF!
    mu_assert_bool(cursor->eof == true);
    mu_assert_bool(cursor->in_session == false);

    // Reuse cursor.
    sky_cursor_set_ptr(cursor, DATA1, DATA1_LENGTH);
    mu_assert_int_equals(cursor->session_event_index, -1);
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int_equals(cursor->session_event_index, 0);
    ASSERT_OBJ_STATE2(cursor->data, 0, "A1", 1000LL, 0LL);
    
    sky_cursor_free(cursor);
    return 0;
}


//--------------------------------------
// Object Iteration
//--------------------------------------

int next_obj(void *_cursor) {
  size_t sz;
  void *ptr = NULL;
  
  sky_cursor *cursor = (sky_cursor*)_cursor;
  if(cursor->context == NULL) {
      ptr = DATA3; sz = DATA3_LENGTH;
  } else if(cursor->context == DATA3) {
      ptr = DATA4; sz = DATA4_LENGTH;
  }
  
  if(ptr != NULL) {
      cursor->context = ptr;
      sky_cursor_set_ptr(cursor, ptr, sz);
      return 1;
  }
  else {
      return 0;
  }
}

int test_sky_cursor_object_iteration() {
    // Setup cursor.
    sky_cursor *cursor = sky_cursor_new(0, 1);
    cursor->next_object_func = next_obj;
    sky_cursor_set_ts_offset(cursor, offsetof(test2_t, ts));
    sky_cursor_set_timestamp_offset(cursor, offsetof(test2_t, timestamp));
    sky_cursor_set_property(cursor, 1, offsetof(test2_t, int_value), sizeof(int32_t), "integer");
    sky_cursor_set_data_sz(cursor, sizeof(test2_t));
    test2_t *obj = (test2_t*)cursor->data;

    // Loop over first object.
    mu_assert_bool(sky_cursor_next_object(cursor));
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int64_equals(obj->int_value, 2LL);
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int64_equals(obj->int_value, 3LL);
    mu_assert_bool(!sky_lua_cursor_next_event(cursor));

    // Loop over second object.
    mu_assert_bool(sky_cursor_next_object(cursor));
    mu_assert_int64_equals(obj->int_value, 0LL);
    mu_assert_bool(sky_lua_cursor_next_event(cursor));
    mu_assert_int64_equals(obj->int_value, 4LL);
    mu_assert_bool(!sky_lua_cursor_next_event(cursor));
    
    // End!
    mu_assert_bool(!sky_cursor_next_object(cursor));

    sky_cursor_free(cursor);
    return 0;
}


//--------------------------------------
// Property Management
//--------------------------------------

int test_sky_cursor_set_integer() {
    size_t sz;
    sky_cursor *cursor = sky_cursor_new(0, 1);
    sky_cursor_set_property(cursor, 1, offsetof(test2_t, int_value), sizeof(int32_t), "integer");
    sky_cursor_set_data_sz(cursor, sizeof(test2_t));
    mu_assert_int_equals(cursor->property_zero_descriptor[1].offset, 8);
    sky_cursor_set_value(cursor, cursor->data, 1, INT_DATA, &sz);
    mu_assert_long_equals(sz, 3L);
    mu_assert_int64_equals(((test2_t*)cursor->data)->int_value, 1000LL);
    sky_cursor_free(cursor);
    return 0;
}

int test_sky_cursor_set_double() {
    size_t sz;
    sky_cursor *cursor = sky_cursor_new(-1, 0);
    sky_cursor_set_property(cursor, -1, offsetof(test2_t, double_value), sizeof(double), "float");
    sky_cursor_set_data_sz(cursor, sizeof(test2_t));
    mu_assert_int_equals(cursor->property_zero_descriptor[-1].offset, 16);
    sky_cursor_set_value(cursor, cursor->data, -1, DOUBLE_DATA, &sz);
    mu_assert_long_equals(sz, 9L);
    mu_assert_bool(fabs(((test2_t*)cursor->data)->double_value - 100.2) < 0.1);
    sky_cursor_free(cursor);
    return 0;
}

int test_sky_cursor_set_boolean() {
    size_t sz;
    sky_cursor *cursor = sky_cursor_new(0, 2);
    sky_cursor_set_property(cursor, 2, offsetof(test2_t, boolean_value), sizeof(bool), "boolean");
    sky_cursor_set_data_sz(cursor, sizeof(test2_t));
    mu_assert_int_equals(cursor->property_zero_descriptor[2].offset, 24);
    sky_cursor_set_value(cursor, cursor->data, 2, BOOLEAN_TRUE_DATA, &sz);
    mu_assert_long_equals(sz, 1L);
    mu_assert_bool(((test2_t*)cursor->data)->boolean_value == true);
    sky_cursor_set_value(cursor, cursor->data, 2, BOOLEAN_FALSE_DATA, &sz);
    mu_assert_long_equals(sz, 1L);
    mu_assert_bool(((test2_t*)cursor->data)->boolean_value == false);
    sky_cursor_free(cursor);
    return 0;
}

int test_sky_cursor_set_string() {
    size_t sz;
    sky_cursor *cursor = sky_cursor_new(0, 1);
    sky_cursor_set_data_sz(cursor, sizeof(test2_t));
    sky_cursor_set_property(cursor, 1, offsetof(test2_t, string_value), sizeof(sky_string), "string");
    mu_assert_int_equals(cursor->property_zero_descriptor[1].offset, 32);
    sky_cursor_set_value(cursor, cursor->data, 1, STRING_DATA, &sz);
    mu_assert_long_equals(sz, 4L);
    mu_assert_int_equals(((test2_t*)cursor->data)->string_value.length, 3);
    mu_assert_bool(((test2_t*)cursor->data)->string_value.data == &STRING_DATA[1]);
    sky_cursor_free(cursor);
    return 0;
}



//==============================================================================
//
// Setup
//
//==============================================================================

int all_tests() {
    mu_run_test(test_sky_cursor_set_data);
    mu_run_test(test_sky_cursor_sessionize);
    mu_run_test(test_sky_cursor_object_iteration);
    
    mu_run_test(test_sky_cursor_set_integer);
    mu_run_test(test_sky_cursor_set_double);
    mu_run_test(test_sky_cursor_set_boolean);
    mu_run_test(test_sky_cursor_set_string);
    return 0;
}

RUN_TESTS()
