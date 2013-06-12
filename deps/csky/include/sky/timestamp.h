#ifndef _sky_timestamp_h
#define _sky_timestamp_h

#include <inttypes.h>

//==============================================================================
//
// Functions
//
//==============================================================================

//--------------------------------------
// Shifting
//--------------------------------------

int64_t sky_timestamp_shift(int64_t value);

int64_t sky_timestamp_unshift(int64_t value);

int64_t sky_timestamp_to_seconds(int64_t value);

#endif

