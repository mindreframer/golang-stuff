package skyd

/*
#include <sky/cursor.h>
#include <leveldb/c.h>
*/
import "C"

import (
	"bytes"
	"unsafe"
)

// NOTE: The callback has to live separate from the execution engine because
// it uses "export".

//export executionEngine_nextObject
func executionEngine_nextObject(cursor unsafe.Pointer) C.int {
	e := (*ExecutionEngine)(((*C.sky_cursor)(cursor)).context)

	// If the iterator is invalid then exit.
	if !e.iterator.Valid() {
		return 0
	}

	// If the key prefix doesn't match then the iterator is done.
	key := e.iterator.Key()
	if !bytes.HasPrefix(key, e.prefix) {
		return 0
	}

	// Set the object data on the cursor.
	value := e.iterator.Value()
	C.sky_cursor_set_ptr(e.cursor, unsafe.Pointer(&value[0]), (C.size_t)(len(value)))

	// Move to the next object.
	e.iterator.Next()

	return 1
}
