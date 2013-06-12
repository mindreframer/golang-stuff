#include <stdio.h>
#include <stdlib.h>
#include "sky/cursor.h"
#include "sky/mem.h"
#include "sky/timestamp.h"
#include "sky/minipack.h"
#include "sky/sky_string.h"
#include "sky/dbg.h"

//==============================================================================
//
// Constants
//
//==============================================================================

#define SKY_PROPERTY_DESCRIPTOR_PADDING  32


//==============================================================================
//
// Macros
//
//==============================================================================

#define badcursordata(MSG, PTR) do {\
    fprintf(stderr, "Cursor pointing at invalid raw event data [" MSG "]: %p->%p\n", cursor->ptr, PTR); \
    memdump(cursor->startptr, (cursor->endptr - cursor->startptr)); \
    cursor->eof = true; \
    return; \
} while(0)


//==============================================================================
//
// Forward Declarations
//
//==============================================================================

//--------------------------------------
// Setters
//--------------------------------------

void sky_set_noop(void *target, void *value, size_t *sz);

void sky_set_string(void *target, void *value, size_t *sz);

void sky_set_int(void *target, void *value, size_t *sz);

void sky_set_double(void *target, void *value, size_t *sz);

void sky_set_boolean(void *target, void *value, size_t *sz);


//--------------------------------------
// Clear Functions
//--------------------------------------

void sky_clear_string(void *target);

void sky_clear_int(void *target);

void sky_clear_double(void *target);

void sky_clear_boolean(void *target);


//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Lifecycle
//--------------------------------------

// Creates a reference to a cursor.
sky_cursor *sky_cursor_new(int32_t min_property_id,
                           int32_t max_property_id)
{
    sky_cursor *cursor = calloc(1, sizeof(sky_cursor));

    // Add one property to account for the zero descriptor.
    min_property_id -= SKY_PROPERTY_DESCRIPTOR_PADDING;
    max_property_id += SKY_PROPERTY_DESCRIPTOR_PADDING;
    int32_t property_count = (max_property_id - min_property_id) + 1;

    // Allocate memory for the descriptors.
    cursor->property_descriptors = calloc(property_count, sizeof(sky_property_descriptor));
    cursor->property_count = property_count;
    cursor->property_zero_descriptor = NULL;
    
    // Initialize all property descriptors to noop.
    int32_t i;
    for(i=0; i<property_count; i++) {
        int64_t property_id = min_property_id + (int64_t)i;
        cursor->property_descriptors[i].property_id = property_id;
        cursor->property_descriptors[i].set_func = sky_set_noop;
        
        // Save a pointer to the descriptor that points to property zero.
        if(property_id == 0) {
            cursor->property_zero_descriptor = &cursor->property_descriptors[i];
        }
    }

    return cursor;
}

// Removes a cursor reference from memory.
void sky_cursor_free(sky_cursor *cursor)
{
    if(cursor) {
        if(cursor->property_descriptors != NULL) free(cursor->property_descriptors);
        cursor->property_zero_descriptor = NULL;
        cursor->property_count = 0;

        if(cursor->data != NULL) free(cursor->data);

        free(cursor);
    }
}


//--------------------------------------
// Data Management
//--------------------------------------

void sky_cursor_set_value(sky_cursor *cursor, void *target,
                          int64_t property_id, void *ptr, size_t *sz)
{
    sky_property_descriptor *property_descriptor = &cursor->property_zero_descriptor[property_id];
    property_descriptor->set_func(target + property_descriptor->offset, ptr, sz);
}


//--------------------------------------
// Descriptor Management
//--------------------------------------

void sky_cursor_set_data_sz(sky_cursor *cursor, uint32_t sz) {
    cursor->data_sz = sz;
    if(cursor->data != NULL) free(cursor->data);
    cursor->data = calloc(1, sz);
}

void sky_cursor_set_timestamp_offset(sky_cursor *cursor, uint32_t offset) {
    cursor->timestamp_descriptor.timestamp_offset = offset;
}

void sky_cursor_set_ts_offset(sky_cursor *cursor, uint32_t offset) {
    cursor->timestamp_descriptor.ts_offset = offset;
}

// Sets the data type and offset for a given property id.
void sky_cursor_set_property(sky_cursor *cursor, int64_t property_id,
                             uint32_t offset, uint32_t sz, const char *data_type)
{
    sky_property_descriptor *property_descriptor = &cursor->property_zero_descriptor[property_id];
    
    // Set the offset and set_func function on the descriptor.
    property_descriptor->offset = offset;
    if(strlen(data_type) == 0) {
        property_descriptor->set_func = sky_set_noop;
        property_descriptor->clear_func = NULL;
    }
    else if(strcmp(data_type, "string") == 0) {
        property_descriptor->set_func = sky_set_string;
        property_descriptor->clear_func = sky_clear_string;
    }
    else if(strcmp(data_type, "factor") == 0 || strcmp(data_type, "integer") == 0) {
        property_descriptor->set_func = sky_set_int;
        property_descriptor->clear_func = sky_clear_int;
    }
    else if(strcmp(data_type, "float") == 0) {
        property_descriptor->set_func = sky_set_double;
        property_descriptor->clear_func = sky_clear_double;
    }
    else if(strcmp(data_type, "boolean") == 0) {
        property_descriptor->set_func = sky_set_boolean;
        property_descriptor->clear_func = sky_clear_boolean;
    }
    else {
        property_descriptor->set_func = sky_set_boolean;
        property_descriptor->clear_func = sky_clear_boolean;
    }
    
    // Resize the action data area.
    if(property_id < 0 && offset+sz > cursor->action_data_sz) {
        cursor->action_data_sz = offset+sz;
    }
}


//--------------------------------------
// Object Iteration
//--------------------------------------

// Moves the cursor to point to the next object.
bool sky_cursor_next_object(sky_cursor *cursor)
{
    return (bool)cursor->next_object_func(cursor);
}


//--------------------------------------
// Event Iteration
//--------------------------------------

void sky_cursor_set_ptr(sky_cursor *cursor, void *ptr, size_t sz)
{
    // Set the start of the path and the length of the data.
    cursor->startptr   = ptr;
    cursor->nextptr    = ptr;
    cursor->endptr     = ptr + sz;
    cursor->ptr        = NULL;
    cursor->in_session = true;
    cursor->last_timestamp      = 0;
    cursor->session_idle_in_sec = 0;
    cursor->session_event_index = -1;
    cursor->eof        = !(ptr != NULL && cursor->startptr < cursor->endptr);
    
    // Clear the data object if set.
    memset(cursor->data, 0, cursor->data_sz);
    
    // The first item is the current state so skip it.
    if(cursor->startptr != NULL && minipack_is_raw(cursor->startptr)) {
        cursor->startptr += minipack_sizeof_elem_and_data(cursor->startptr);
        cursor->nextptr = cursor->startptr;
    }
}

void sky_cursor_next_event(sky_cursor *cursor)
{
    // Ignore any calls when the cursor is out of session or EOF.
    if(cursor->eof || !cursor->in_session) {
        return;
    }

    // Move the pointer to the next position.
    void *prevptr = cursor->ptr;
    cursor->ptr = cursor->nextptr;
    void *ptr = cursor->ptr;

    // If pointer is beyond the last event then set eof.
    if(cursor->ptr >= cursor->endptr) {
        cursor->eof        = true;
        cursor->in_session = false;
        cursor->ptr        = NULL;
        cursor->startptr   = NULL;
        cursor->nextptr    = NULL;
        cursor->endptr     = NULL;
    }
    // Otherwise update the event object with data.
    else {
        sky_event_flag_t flag = *((sky_event_flag_t*)ptr);
        
        // If flag isn't correct then report and exit.
        if(flag != EVENT_FLAG) badcursordata("eflag", ptr);
        ptr += sizeof(sky_event_flag_t);
        
        // Read timestamp.
        size_t sz;
        int64_t ts = minipack_unpack_int(ptr, &sz);
        if(sz == 0) badcursordata("timestamp", ptr);
        uint32_t timestamp = sky_timestamp_to_seconds(ts);
        ptr += sz;

        // Check for session boundry. This only applies if this is not the
        // first event in the session and a session idle time has been set.
        if(cursor->last_timestamp > 0 && cursor->session_idle_in_sec > 0) {
            // If the elapsed time is greater than the idle time then rewind
            // back to the event we started on at the beginning of the function
            // and mark the cursor as being "out of session".
            if(timestamp - cursor->last_timestamp >= cursor->session_idle_in_sec) {
                cursor->ptr = prevptr;
                cursor->in_session = false;
            }
        }
        cursor->last_timestamp = timestamp;

        // Only process the event if we're still in session.
        if(cursor->in_session) {
            cursor->session_event_index++;
            
            // Set timestamp.
            int64_t *data_ts = (int64_t*)(cursor->data + cursor->timestamp_descriptor.ts_offset);
            uint32_t *data_timestamp = (uint32_t*)(cursor->data + cursor->timestamp_descriptor.timestamp_offset);
            *data_ts = ts;
            *data_timestamp = timestamp;
            
            // Clear old action data.
            if(cursor->action_data_sz > 0) {
              memset(cursor->data, 0, cursor->action_data_sz);
            }

            // Read msgpack map!
            uint32_t count = minipack_unpack_map(ptr, &sz);
            if(sz == 0) {
              minipack_unpack_nil(ptr, &sz);
              if(sz == 0) {
                badcursordata("datamap", ptr);
              }
            }
            ptr += sz;

            // Loop over key/value pairs.
            uint32_t i;
            for(i=0; i<count; i++) {
                // Read property id (key).
                int64_t property_id = minipack_unpack_int(ptr, &sz);
                if(sz == 0) badcursordata("key", ptr);
                ptr += sz;

                // Read property value and set it on the data object.
                sky_cursor_set_value(cursor, cursor->data, property_id, ptr, &sz);
                if(sz == 0) {
                  debug("[invalid read, skipping]");
                  sz = minipack_sizeof_elem_and_data(ptr);
                }
                ptr += sz;
            }

            cursor->nextptr = ptr;
        }
    }
}

bool sky_lua_cursor_next_event(sky_cursor *cursor)
{
    sky_cursor_next_event(cursor);
    return (!cursor->eof && cursor->in_session);
}

bool sky_cursor_eof(sky_cursor *cursor)
{
    return cursor->eof;
}

bool sky_cursor_eos(sky_cursor *cursor)
{
    return !cursor->in_session;
}

void sky_cursor_set_session_idle(sky_cursor *cursor, uint32_t seconds)
{
    // Save the idle value.
    cursor->session_idle_in_sec = seconds;

    // If the value is non-zero then start sessionizing the cursor.
    cursor->in_session = (seconds > 0 ? false : !cursor->eof);
}

void sky_cursor_next_session(sky_cursor *cursor)
{
    // Set a flag to allow the cursor to continue iterating unless EOF is set.
    if(!cursor->in_session) {
        cursor->session_event_index = -1;
        cursor->in_session = !cursor->eof;
    }
}

bool sky_lua_cursor_next_session(sky_cursor *cursor)
{
    sky_cursor_next_session(cursor);
    return !cursor->eof;
}



//--------------------------------------
// Setters
//--------------------------------------

void sky_set_noop(void *target, void *value, size_t *sz)
{
    ((void)(target));
    *sz = minipack_sizeof_elem_and_data(value);
}

void sky_set_string(void *target, void *value, size_t *sz)
{
    size_t _sz;
    sky_string *string = (sky_string*)target;
    string->length = minipack_unpack_raw(value, &_sz);
    string->data = (_sz > 0 ? value + _sz : NULL);
    *sz = _sz + string->length;
}

void sky_set_int(void *target, void *value, size_t *sz)
{
    *((int32_t*)target) = (int32_t)minipack_unpack_int(value, sz);
}

void sky_set_double(void *target, void *value, size_t *sz)
{
    *((double*)target) = minipack_unpack_double(value, sz);
}

void sky_set_boolean(void *target, void *value, size_t *sz)
{
    *((bool*)target) = minipack_unpack_bool(value, sz);
}


//--------------------------------------
// Clear Functions
//--------------------------------------

void sky_clear_string(void *target)
{
    sky_string *string = (sky_string*)target;
    string->length = 0;
    string->data = NULL;
}

void sky_clear_int(void *target)
{
    *((int32_t*)target) = 0;
}

void sky_clear_double(void *target)
{
    *((double*)target) = 0;
}

void sky_clear_boolean(void *target)
{
    *((bool*)target) = false;
}

