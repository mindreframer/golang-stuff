#ifndef _sky_cursor_h
#define _sky_cursor_h

#include <stdlib.h>
#include <inttypes.h>
#include <stdbool.h>


//==============================================================================
//
// Constants
//
//==============================================================================

#define sky_event_flag_t uint8_t
#define EVENT_FLAG       0x92


//==============================================================================
//
// Typedefs
//
//==============================================================================

typedef struct sky_cursor sky_cursor;

typedef int (*sky_cursor_next_object_func)(void *cursor);

typedef void (*sky_property_descriptor_set_func)(void *target, void *value, size_t *sz);
typedef void (*sky_property_descriptor_clear_func)(void *target);

typedef struct { uint16_t ts_offset; uint16_t timestamp_offset;} sky_timestamp_descriptor;

typedef struct {
    int64_t property_id;
    uint16_t offset;
    sky_property_descriptor_set_func set_func;
    sky_property_descriptor_clear_func clear_func;
} sky_property_descriptor;

struct sky_cursor {
    void *data;
    uint32_t data_sz;
    uint32_t action_data_sz;

    int32_t session_event_index;
    void *startptr;
    void *nextptr;
    void *endptr;
    void *ptr;
    bool eof;
    bool in_session;
    uint32_t last_timestamp;
    uint32_t session_idle_in_sec;

    sky_timestamp_descriptor timestamp_descriptor;
    sky_property_descriptor *property_descriptors;
    sky_property_descriptor *property_zero_descriptor;
    uint32_t property_count;

    void *context;
    sky_cursor_next_object_func next_object_func;
};


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

sky_cursor *sky_cursor_new(int32_t min_property_id, int32_t max_property_id);

void sky_cursor_free(sky_cursor *cursor);


//--------------------------------------
// Data Management
//--------------------------------------

void sky_cursor_set_value(sky_cursor *cursor,
  void *target, int64_t property_id, void *ptr, size_t *sz);


//--------------------------------------
// Descriptor Management
//--------------------------------------

void sky_cursor_set_data_sz(sky_cursor *cursor, uint32_t sz);

void sky_cursor_set_timestamp_offset(sky_cursor *cursor, uint32_t offset);

void sky_cursor_set_ts_offset(sky_cursor *cursor, uint32_t offset);

void sky_cursor_set_property(sky_cursor *cursor,
  int64_t property_id, uint32_t offset, uint32_t sz, const char *data_type);

//--------------------------------------
// Object Iteration
//--------------------------------------

bool sky_cursor_next_object(sky_cursor *cursor);


//--------------------------------------
// Event Iteration
//--------------------------------------

void sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz);

void sky_cursor_next_event(sky_cursor *cursor);

bool sky_lua_cursor_next_event(sky_cursor *cursor);

bool sky_cursor_eof(sky_cursor *cursor);

bool sky_cursor_eos(sky_cursor *cursor);

void sky_cursor_set_session_idle(sky_cursor *cursor, uint32_t seconds);

void sky_cursor_next_session(sky_cursor *cursor);

bool sky_lua_cursor_next_session(sky_cursor *cursor);

void sky_cursor_clear_data(sky_cursor *cursor);

#endif