package skyd

/*
#cgo LDFLAGS: -lcsky -lluajit-5.1 -lleveldb
#include <stdlib.h>
#include <sky/cursor.h>
#include <luajit-2.0/lua.h>
#include <luajit-2.0/lualib.h>
#include <luajit-2.0/lauxlib.h>

int mp_pack(lua_State *L);
int mp_unpack(lua_State *L);

int executionEngine_nextObject(void *cursor);

int executionEngine_c_next_object(void *cursor) {
	return (bool)executionEngine_nextObject(cursor);
}

void executionEngine_setNextObjectFunc(void *cursor) {
	((sky_cursor*)cursor)->next_object_func = executionEngine_c_next_object;
}

*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jmhodges/levigo"
	"github.com/ugorji/go-msgpack"
	"regexp"
	"sort"
	"text/template"
	"unsafe"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// An ExecutionEngine is used to iterate over a series of objects.
type ExecutionEngine struct {
	tableName    string
	iterator     *levigo.Iterator
	cursor       *C.sky_cursor
	prefix       []byte
	state        *C.lua_State
	header       string
	source       string
	fullSource   string
	propertyFile *PropertyFile
	propertyRefs []*Property

	cprefix    unsafe.Pointer
	cprefix_sz C.size_t
}

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

func NewExecutionEngine(table *Table, source string) (*ExecutionEngine, error) {
	if table == nil {
		return nil, errors.New("skyd.ExecutionEngine: Table required")
	}
	propertyFile := table.propertyFile
	if propertyFile == nil {
		return nil, errors.New("skyd.ExecutionEngine: Property file required")
	}

	// Find a list of all references properties.
	propertyRefs, err := extractPropertyReferences(propertyFile, source)
	if err != nil {
		return nil, err
	}

	// Determine table prefix.
	prefix, err := TablePrefix(table.Name)
	if err != nil {
		return nil, err
	}

	// Create the engine.
	e := &ExecutionEngine{
		tableName:    table.Name,
		prefix:       prefix,
		propertyFile: propertyFile,
		source:       source,
		propertyRefs: propertyRefs,
	}

	// Initialize the engine.
	err = e.init()
	if err != nil {
		fmt.Printf("%s\n\n", e.FullAnnotatedSource())
		e.Destroy()
		return nil, err
	}

	return e, nil
}

//------------------------------------------------------------------------------
//
// Properties
//
//------------------------------------------------------------------------------

// Retrieves the source for the engine.
func (e *ExecutionEngine) Source() string {
	return e.source
}

// Retrieves the generated header for the engine.
func (e *ExecutionEngine) Header() string {
	return e.header
}

// Retrieves the full source sent to the Lua compiler.
func (e *ExecutionEngine) FullSource() string {
	return e.fullSource
}

// Retrieves the full annotated source with line numbers.
func (e *ExecutionEngine) FullAnnotatedSource() string {
	lineNumber := 1
	r, _ := regexp.Compile(`\n`)
	return "00001 " + r.ReplaceAllStringFunc(e.fullSource, func(str string) string {
		lineNumber += 1
		return fmt.Sprintf("%s%05d ", str, lineNumber)
	})
}

// Sets the iterator to use.
func (e *ExecutionEngine) SetIterator(iterator *levigo.Iterator) error {
	// Close the old iterator.
	if e.iterator != nil {
		e.iterator.Close()
	}

	// Attach the new iterator.
	e.iterator = iterator
	if e.iterator != nil {
		e.iterator.Seek(e.prefix)
	}

	return nil
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Lifecycle
//--------------------------------------

// Initializes the Lua context and compiles the source code.
func (e *ExecutionEngine) init() error {
	if e.state != nil {
		return nil
	}

	// Initialize the state and open the libraries.
	e.state = C.luaL_newstate()
	if e.state == nil {
		return errors.New("Unable to initialize Lua context.")
	}
	C.luaL_openlibs(e.state)

	// Generate the header file.
	err := e.generateHeader()
	if err != nil {
		e.Destroy()
		return err
	}

	// Compile the script.
	e.fullSource = fmt.Sprintf("%v\n%v", e.header, e.source)
	source := C.CString(e.fullSource)
	defer C.free(unsafe.Pointer(source))
	ret := C.luaL_loadstring(e.state, source)
	if ret != 0 {
		defer e.Destroy()
		errstring := C.GoString(C.lua_tolstring(e.state, -1, nil))
		return fmt.Errorf("skyd.ExecutionEngine: Syntax Error: %v", errstring)
	}

	// Run script once to initialize.
	ret = C.lua_pcall(e.state, 0, 0, 0)
	if ret != 0 {
		defer e.Destroy()
		errstring := C.GoString(C.lua_tolstring(e.state, -1, nil))
		return fmt.Errorf("skyd.ExecutionEngine: Init Error: %v", errstring)
	}

	// Setup cursor.
	err = e.initCursor()
	if err != nil {
		e.Destroy()
		return err
	}

	return nil
}

// Initializes the cursor used by the script.
func (e *ExecutionEngine) initCursor() error {
	// Create the cursor.
	minPropertyId, maxPropertyId := e.propertyFile.NextIdentifiers()
	e.cursor = C.sky_cursor_new((C.int32_t)(minPropertyId), (C.int32_t)(maxPropertyId))
	e.cursor.context = unsafe.Pointer(e)
	C.executionEngine_setNextObjectFunc(unsafe.Pointer(e.cursor))

	// Initialize the cursor from within Lua.
	functionName := C.CString("sky_init_cursor")
	defer C.free(unsafe.Pointer(functionName))

	C.lua_getfield(e.state, -10002, functionName)
	C.lua_pushlightuserdata(e.state, unsafe.Pointer(e.cursor))
	//fmt.Printf("%s\n\n", e.FullAnnotatedSource())
	rc := C.lua_pcall(e.state, 1, 0, 0)
	if rc != 0 {
		luaErrString := C.GoString(C.lua_tolstring(e.state, -1, nil))
		return fmt.Errorf("Unable to init cursor: %s", luaErrString)
	}

	return nil
}

// Closes the lua context.
func (e *ExecutionEngine) Destroy() {
	if e.state != nil {
		C.lua_close(e.state)
		e.state = nil
	}
	if e.iterator != nil {
		e.SetIterator(nil)
	}
}

//--------------------------------------
// Execution
//--------------------------------------

// Executes an aggregation over the iterator.
func (e *ExecutionEngine) Aggregate() (interface{}, error) {
	functionName := C.CString("sky_aggregate")
	defer C.free(unsafe.Pointer(functionName))

	C.lua_getfield(e.state, -10002, functionName)
	C.lua_pushlightuserdata(e.state, unsafe.Pointer(e.cursor))
	rc := C.lua_pcall(e.state, 1, 1, 0)
	if rc != 0 {
		luaErrString := C.GoString(C.lua_tolstring(e.state, -1, nil))
		fmt.Println(e.FullAnnotatedSource())
		return nil, fmt.Errorf("skyd.ExecutionEngine: Unable to aggregate: %s", luaErrString)
	}

	return e.decodeResult()
}

// Executes an merge over the iterator.
func (e *ExecutionEngine) Merge(results interface{}, data interface{}) (interface{}, error) {
	functionName := C.CString("sky_merge")
	defer C.free(unsafe.Pointer(functionName))

	C.lua_getfield(e.state, -10002, functionName)
	err := e.encodeArgument(results)
	if err != nil {
		return results, err
	}
	err = e.encodeArgument(data)
	if err != nil {
		return results, err
	}
	rc := C.lua_pcall(e.state, 2, 1, 0)
	if rc != 0 {
		luaErrString := C.GoString(C.lua_tolstring(e.state, -1, nil))
		fmt.Println(e.FullAnnotatedSource())
		return results, fmt.Errorf("skyd.ExecutionEngine: Unable to merge: %s", luaErrString)
	}

	return e.decodeResult()
}

// Encodes a Go object into Msgpack and adds it to the function arguments.
func (e *ExecutionEngine) encodeArgument(value interface{}) error {
	// Encode Go object into msgpack.
	buffer := new(bytes.Buffer)
	encoder := msgpack.NewEncoder(buffer)
	err := encoder.Encode(value)
	if err != nil {
		return err
	}

	// Push the msgpack data onto the Lua stack.
	data := buffer.String()
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))
	C.lua_pushlstring(e.state, cdata, (C.size_t)(len(data)))

	// Convert the argument from msgpack into Lua.
	rc := C.mp_unpack(e.state)
	if rc != 1 {
		return errors.New("skyd.ExecutionEngine: Unable to msgpack encode Lua argument")
	}
	C.lua_remove(e.state, -2)

	return nil
}

// Decodes the result from a function into a Go object.
func (e *ExecutionEngine) decodeResult() (interface{}, error) {
	// Encode Lua object into msgpack.
	rc := C.mp_pack(e.state)
	if rc != 1 {
		return nil, errors.New("skyd.ExecutionEngine: Unable to msgpack decode Lua result")
	}
	sz := C.size_t(0)
	ptr := C.lua_tolstring(e.state, -1, (*C.size_t)(&sz))
	str := C.GoStringN(ptr, (C.int)(sz))
	C.lua_settop(e.state, -(1)-1) // lua_pop()

	// Decode msgpack into a Go object.
	var ret interface{}
	decoder := msgpack.NewDecoder(bytes.NewBufferString(str), nil)
	err := decoder.Decode(&ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

//--------------------------------------
// Codegen
//--------------------------------------

// Generates the header for the script based on a source string.
func (e *ExecutionEngine) generateHeader() error {
	// Parse the header template.
	t := template.New("header.lua")
	t.Funcs(template.FuncMap{"structdef": propertyStructDef, "metatypedef": metatypeFunctionDef, "initdescriptor": initDescriptorDef})
	_, err := t.Parse(LuaHeader)
	if err != nil {
		return err
	}

	// Generate the template from the property references.
	var buffer bytes.Buffer
	err = t.Execute(&buffer, e.propertyRefs)
	if err != nil {
		return err
	}

	// Assign header
	e.header = buffer.String()

	return nil
}

// Extracts the property references from the source string.
func extractPropertyReferences(propertyFile *PropertyFile, source string) ([]*Property, error) {
	// Create a list of properties.
	properties := make([]*Property, 0)
	lookup := make(map[int64]*Property)

	// Find all the event property references in the script.
	r, err := regexp.Compile(`\bevent(?:\.|:)(\w+)`)
	if err != nil {
		return nil, err
	}
	for _, match := range r.FindAllStringSubmatch(source, -1) {
		name := match[1]
		property := propertyFile.GetPropertyByName(name)
		if property == nil {
			return nil, fmt.Errorf("Property not found: '%v'", name)
		}
		if lookup[property.Id] == nil {
			properties = append(properties, property)
			lookup[property.Id] = property
		}
	}
	sort.Sort(PropertyList(properties))

	return properties, nil
}

func propertyStructDef(args ...interface{}) string {
	if property, ok := args[0].(*Property); ok {
		return fmt.Sprintf("%v _%v;", getPropertyCType(property), property.Name)
	}
	return ""
}

func metatypeFunctionDef(args ...interface{}) string {
	if property, ok := args[0].(*Property); ok {
		switch property.DataType {
		case StringDataType:
			return fmt.Sprintf("%v = function(event) return ffi.string(event._%v.data, event._%v.length) end,", property.Name, property.Name, property.Name)
		default:
			return fmt.Sprintf("%v = function(event) return event._%v end,", property.Name, property.Name)
		}
	}
	return ""
}

func initDescriptorDef(args ...interface{}) string {
	if property, ok := args[0].(*Property); ok {
		return fmt.Sprintf("cursor:set_property(%d, ffi.offsetof('sky_lua_event_t', '_%s'), ffi.sizeof('%s'), '%s')", property.Id, property.Name, getPropertyCType(property), property.DataType)
	}
	return ""
}

func getPropertyCType(property *Property) string {
	switch property.DataType {
	case StringDataType:
		return "sky_string_t"
	case FactorDataType, IntegerDataType:
		return "int32_t"
	case FloatDataType:
		return "double"
	case BooleanDataType:
		return "bool"
	default:
		panic(fmt.Sprintf("skyd.ExecutionEngine: Invalid data type: %v", property.DataType))
	}
	return ""
}
