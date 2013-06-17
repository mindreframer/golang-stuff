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

const (
	// Broadcast Channel Size
	broadcastChannelSize = 16
)

// The Hub manages all active connection, but should only be used directly
// if broadcasting of data or an event to all connections is desired.
// The Hub should not be instanced directly. Use GetHub to get the active hub
// for broadcasting messages.
type Hub struct {
	// Registered connections.
	connections map[*Connection]bool

	// Inbound messages from the connections.
	broadcast chan *message

	// Register requests from the connections.
	register chan *Connection

	// Unregister requests from connections.
	unregister chan *Connection

	// Flag to determine if running or not
	isRunning bool
}

// Remove the specified connection from the hub and drop the socket.
func (hub *Hub) remove(conn *Connection) {
	delete(hub.connections, conn)
	close(conn.send)
}

// If the hub is not running, start it in a different goroutine.
func (hub *Hub) run() {
	if hub.isRunning != true { // Should be safe, because only called from NewRouter and therefore a single thread.
		hub.isRunning = true
		go func() {
			for {
				select {
				// Register new connection
				case conn := <-hub.register:
					hub.connections[conn] = true
				// Unregister dropped connection
				case conn := <-hub.unregister:
					if _, ok := hub.connections[conn]; ok {
						hub.remove(conn)
					}
				// Broadcast
				case message := <-hub.broadcast:
					for conn := range hub.connections {
						select {
						case conn.send <- message:
						default:
							hub.remove(conn)
						}
					}
				}
			}
		}()
	}
}

// Create the hub instance.
var hub = Hub{
	broadcast:   make(chan *message, broadcastChannelSize),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	connections: make(map[*Connection]bool),
	isRunning:   false,
}

// GetHub retrieves and returns pointer to golem's active hub.
func GetHub() *Hub {
	return &hub
}

// Broadcast emits an event with data to ALL active connections.
func (hub *Hub) Broadcast(event string, data interface{}) {
	hub.broadcast <- &message{
		event: event,
		data:  data,
	}
}
