package skyd

const LuaHeader = `
-- SKY GENERATED CODE BEGIN --
local ffi = require('ffi')
ffi.cdef([[
typedef struct sky_string_t { int32_t length; char *data; } sky_string_t;
typedef struct {
  {{range .}}{{structdef .}}
  {{end}}
  int64_t ts;
  uint32_t timestamp;
} sky_lua_event_t;
typedef struct sky_cursor_t { sky_lua_event_t *event; int32_t session_event_index; } sky_cursor_t;

int sky_cursor_set_data_sz(sky_cursor_t *cursor, uint32_t sz);
int sky_cursor_set_timestamp_offset(sky_cursor_t *cursor, uint32_t offset);
int sky_cursor_set_ts_offset(sky_cursor_t *cursor, uint32_t offset);
int sky_cursor_set_property(sky_cursor_t *cursor, int64_t property_id, uint32_t offset, uint32_t sz, const char *data_type);

bool sky_cursor_has_next_object(sky_cursor_t *);
bool sky_cursor_next_object(sky_cursor_t *);
bool sky_cursor_eof(sky_cursor_t *);
bool sky_cursor_eos(sky_cursor_t *);
bool sky_lua_cursor_next_event(sky_cursor_t *);
bool sky_lua_cursor_next_session(sky_cursor_t *);
bool sky_cursor_set_session_idle(sky_cursor_t *, uint32_t);
]])
ffi.metatype('sky_cursor_t', {
  __index = {
    set_data_sz = function(cursor, sz) return ffi.C.sky_cursor_set_data_sz(cursor, sz) end,
    set_timestamp_offset = function(cursor, offset) return ffi.C.sky_cursor_set_timestamp_offset(cursor, offset) end,
    set_ts_offset = function(cursor, offset) return ffi.C.sky_cursor_set_ts_offset(cursor, offset) end,
    set_action_id_offset = function(cursor, offset) return ffi.C.sky_cursor_set_action_id_offset(cursor, offset) end,
    set_property = function(cursor, property_id, offset, sz, data_type) return ffi.C.sky_cursor_set_property(cursor, property_id, offset, sz, data_type) end,

    hasNextObject = function(cursor) return ffi.C.sky_cursor_has_next_object(cursor) end,
    nextObject = function(cursor) return ffi.C.sky_cursor_next_object(cursor) end,
    eof = function(cursor) return ffi.C.sky_cursor_eof(cursor) end,
    eos = function(cursor) return ffi.C.sky_cursor_eos(cursor) end,
    next = function(cursor) return ffi.C.sky_lua_cursor_next_event(cursor) end,
    next_session = function(cursor) return ffi.C.sky_lua_cursor_next_session(cursor) end,
    set_session_idle = function(cursor, seconds) return ffi.C.sky_cursor_set_session_idle(cursor, seconds) end,
  }
})
ffi.metatype('sky_lua_event_t', {
  __index = {
  {{range .}}{{metatypedef .}}
  {{end}}
  }
})

function sky_init_cursor(_cursor)
  cursor = ffi.cast('sky_cursor_t*', _cursor)
  {{range .}}{{initdescriptor .}}
  {{end}}
  cursor:set_timestamp_offset(ffi.offsetof('sky_lua_event_t', 'timestamp'))
  cursor:set_ts_offset(ffi.offsetof('sky_lua_event_t', 'ts'))
  cursor:set_data_sz(ffi.sizeof('sky_lua_event_t'))
end

function sky_aggregate(_cursor)
  cursor = ffi.cast('sky_cursor_t*', _cursor)
  data = {}
  while cursor:nextObject() do
    aggregate(cursor, data)
  end
  return data
end

-- The wrapper for the merge.
function sky_merge(results, data)
  if data ~= nil then
    merge(results, data)
  end
  return results
end
-- SKY GENERATED CODE END --
`
