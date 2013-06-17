/*

   Copyright 2013 Niklas Voss

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package golem

import (
	"errors"
	"github.com/garyburd/go-websocket/websocket"
	"log"
	"net/http"
	"reflect"
)

// Router handles multiplexing of incoming messenges by typenames/events.
// Initially a router uses heartbeats and the default protocol.
type Router struct {
	// Map of callbacks for event types.
	callbacks map[string]func(*Connection, interface{})
	// Protocol extensions
	extensions map[reflect.Type]reflect.Value
	// Function being called if connection is closed.
	closeFunc func(*Connection)
	// Function verifying handshake.
	handshakeFunc func(http.ResponseWriter, *http.Request) bool
	// Active protocol
	protocol Protocol
	// Flag to enable or disable heartbeats
	useHeartbeats bool
	//
	connExtensionConstructor reflect.Value
}

// NewRouter intialises a new instance and returns the pointer.
func NewRouter() *Router {
	// Tries to run hub, if already running nothing will happen.
	hub.run()
	// Returns pointer to instance.
	return &Router{
		callbacks:                make(map[string]func(*Connection, interface{})),
		extensions:               make(map[reflect.Type]reflect.Value),
		closeFunc:                func(*Connection) {},                                          // Empty placeholder close function.
		handshakeFunc:            func(http.ResponseWriter, *http.Request) bool { return true }, // Handshake always allowed.
		protocol:                 initialProtocol,
		useHeartbeats:            true,
		connExtensionConstructor: defaultConnectionExtension,
	}
}

// Handler creates a handler function for this router, that can be used with the
// http-package to handle WebSocket-Connections.
func (router *Router) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if method used was GET.
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		// Disallow cross-origin connections.
		if r.Header.Get("Origin") != "http://"+r.Host {
			http.Error(w, "Origin not allowed", 403)
			return
		}
		// Check if handshake callback verifies upgrade.
		if !router.handshakeFunc(w, r) {
			http.Error(w, "Authorization failed", 403)
			return
		}
		// Upgrade websocket connection.
		socket, err := websocket.Upgrade(w, r.Header, nil, 1024, 1024)
		// Check if handshake was successful
		if _, ok := err.(websocket.HandshakeError); ok {
			http.Error(w, "Not a websocket handshake", 400)
			return
		} else if err != nil {
			log.Println(err)
			return
		}

		// Create the connection.
		conn := newConnection(socket, router)
		//
		if router.connExtensionConstructor.IsValid() {
			conn.extend(router.connExtensionConstructor.Call([]reflect.Value{reflect.ValueOf(conn)})[0].Interface())
		}
		// And start reading and writing routines.
		conn.run()
	}
}

// The On-function adds callbacks by name of the event, that should be handled.
// For type T the callback would be of type:
//     func (*golem.Connection, *T)
// Type T can be any type. By default golem tries to unmarshal json into the
// specified type. If a custom protocol is used, it will be used instead to process the data.
// If type T is registered to use a protocol extension, it will be used instead.
// If type T is interface{} the interstage data of the active protocol will be directly forwarded!
// (Note: the golem wiki has a whole page about this function)
func (router *Router) On(name string, callback interface{}) {

	callbackValue := reflect.ValueOf(callback)
	callbackType := reflect.TypeOf(callback)
	if router.connExtensionConstructor.IsValid() {
		extType := router.connExtensionConstructor.Type().Out(0)
		if callbackType.In(0) == extType {
			// EXTENSION TYPE

			// NO DATA
			if callbackType.NumIn() == 1 {
				router.callbacks[name] = func(conn *Connection, data interface{}) {
					args := []reflect.Value{reflect.ValueOf(conn.extension)}
					callbackValue.Call(args)
				}
				return
			}

			// INTERFACE
			if callbackType.In(1).Kind() == reflect.Interface {
				router.callbacks[name] = func(conn *Connection, data interface{}) {
					args := []reflect.Value{reflect.ValueOf(conn.extension), reflect.ValueOf(data)}
					callbackValue.Call(args)
				}
				return
			}

			// PROTOCOL EXTENSION
			if parser, ok := router.extensions[callbackType.In(1)]; ok {
				router.callbacks[name] = func(conn *Connection, data interface{}) {
					if result := parser.Call([]reflect.Value{reflect.ValueOf(data)}); result[1].Bool() {
						args := []reflect.Value{reflect.ValueOf(conn.extension), result[0]}
						callbackValue.Call(args)
					}
				}
				return
			}

			// PROTOCOL
			callbackDataElem := callbackType.In(1).Elem()
			router.callbacks[name] = func(conn *Connection, data interface{}) {
				result := reflect.New(callbackDataElem)

				err := router.protocol.Unmarshal(data, result.Interface())
				if err == nil {
					args := []reflect.Value{reflect.ValueOf(conn.extension), result}
					callbackValue.Call(args)
				} else {
					// TODO: Proper debug output!
				}
			}
			return
		}
	} else {
		// DEFAULT TYPE

		// NO DATA
		if reflect.TypeOf(callback).NumIn() == 1 {
			router.callbacks[name] = func(conn *Connection, data interface{}) {
				callback.(func(*Connection))(conn)
			}
			return
		}

		// INTERFACE
		if cb, ok := callback.(func(*Connection, interface{})); ok {
			router.callbacks[name] = cb
			return
		}

		// PROTOCOL EXTENSION
		if parser, ok := router.extensions[callbackType.In(1)]; ok {
			router.callbacks[name] = func(conn *Connection, data interface{}) {
				if result := parser.Call([]reflect.Value{reflect.ValueOf(data)}); result[1].Bool() {
					args := []reflect.Value{reflect.ValueOf(conn), result[0]}
					callbackValue.Call(args)
				}
			}
			return
		}

		// PROTOCOL
		callbackDataElem := callbackType.In(1).Elem()
		router.callbacks[name] = func(conn *Connection, data interface{}) {
			result := reflect.New(callbackDataElem)

			err := router.protocol.Unmarshal(data, result.Interface())
			if err == nil {
				args := []reflect.Value{reflect.ValueOf(conn), result}
				callbackValue.Call(args)
			} else {
				// TODO: Proper debug output!
			}
		}
		return
	}
}

// Unpacks incoming data and forwards it to callback.
func (router *Router) processMessage(conn *Connection, in []byte) {
	if name, data, err := router.protocol.Unpack(in); err == nil {
		if callback, ok := router.callbacks[name]; ok {
			callback(conn, data)
		}
	} // TODO: else error logging?

	defer recover()
}

// OnClose sets the callback for connection closes.
func (router *Router) OnClose(callback func(*Connection)) {
	router.closeFunc = callback
}

// OnHandshake sets the callback for handshake verfication.
// If the handshake function returns false the request will not be upgraded.
func (router *Router) OnHandshake(callback func(http.ResponseWriter, *http.Request) bool) {
	router.handshakeFunc = callback
}

// The AddProtocolExtension-function allows adding of custom parsers for custom types. For any Type T
// the parser function would look like this:
//      func (interface{}) (T, bool)
// The interface's type is depending on the interstage product of the active protocol, by default
// for the JSON-based protocol it is []byte and therefore the function could be simplified to:
//      func ([]byte) (T, bool)
// Or in general if P is the interstage product:
//      func (*P) (T, bool)
// The boolean return value is necessary to verify if parsing was successful.
// All On-handling function accepting T as input data will now automatically use the custom
// extension. For an example see the example_data.go file in the example repository.
func (router *Router) AddProtocolExtension(extensionFunc interface{}) error {
	extensionValue := reflect.ValueOf(extensionFunc)
	extensionType := extensionValue.Type()

	if extensionType.NumIn() != 1 {
		return errors.New("Cannot add function(" + extensionType.String() + ") as parser: To many arguments!")
	}
	if extensionType.NumOut() != 2 {
		return errors.New("Cannot add function(" + extensionType.String() + ") as parser: Wrong number of return values!")
	}
	if extensionType.Out(1).Kind() != reflect.Bool {
		return errors.New("Cannot add function(" + extensionType.String() + ") as parser: Second return value is not Bool!")
	}

	router.extensions[extensionType.Out(0)] = extensionValue

	return nil
}

// SetConnectionExtension sets the extension for this router. A connection extension is an extended connection structure, that afterwards can be
// used by On-handlers as well (use-cases: add persistent storage to connection, additional methods et cetera).
// The SetConnectionExtension function takes the constructor of the custom format to be able to use and create it on
// connection to the router.
// For type E the constructor needs to fulfil the following requirements:
//     func NewE(conn *Connection) *E
// Afterwards On-handler can us this extended type:
//     router.On(func (extendedConn E, data Datatype) { ... })
// For an example have a look at the example repository and have a look at the 'example_connection_extension.go'.
func (router *Router) SetConnectionExtension(constructor interface{}) {
	router.connExtensionConstructor = reflect.ValueOf(constructor)
}

// SetProtocol sets the protocol of the router to the supplied implementation of the Protocol interface.
func (router *Router) SetProtocol(protocol Protocol) {
	router.protocol = protocol
}

//

// SetHeartbeat activates or deactivates the heartbeat depending on the flag parameter. By default heartbeats are activated.
func (router *Router) SetHeartbeat(flag bool) {
	router.useHeartbeats = flag
}
