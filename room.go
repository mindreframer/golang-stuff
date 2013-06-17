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
	roomSendChannelSize = 32
)

// Rooms are groups of connections. A room provides methods to communicate with all
// members of the group simultaneously.
type Room struct {
	// Map of member connections
	members map[*Connection]bool
	// Stop channel
	stop chan bool
	// Join request
	join chan *Connection
	// Leave request
	leave chan *Connection
	// Broadcast to room members
	send chan *message
}

// Creates and initialised a room and returns pointer to it.
func NewRoom() *Room {
	r := Room{
		members: make(map[*Connection]bool),
		stop:    make(chan bool),
		join:    make(chan *Connection),
		leave:   make(chan *Connection),
		send:    make(chan *message, roomSendChannelSize),
	}
	// Run the message loop
	go r.run()
	// Return pointer
	return &r
}

// Starts the message loop of this room, should only be run once and in a different routine.
func (r *Room) run() {
	for {
		select {
		// Join
		case conn := <-r.join:
			r.members[conn] = true
		// Leave
		case conn := <-r.leave:
			if _, ok := r.members[conn]; ok { // If member exists, delete it
				delete(r.members, conn)
			}
		// Send
		case message := <-r.send:
			for conn := range r.members { // For every connection try to send
				select {
				case conn.send <- message:
				default: // If sending failed, delete member
					delete(r.members, conn)
				}
			}
		// Stop
		case <-r.stop:
			return
		}
	}
}

// Stops and shutsdown the room. After calling Stop the room can be safely deleted.
func (r *Room) Stop() {
	r.stop <- true
}

// Join adds the provided connection to the room.
func (r *Room) Join(conn *Connection) {
	r.join <- conn
}

// Leave removes the connection from the room, if it previously was member of the room.
func (r *Room) Leave(conn *Connection) {
	r.leave <- conn
}

// Emits message event to all members of the room.
func (r *Room) Emit(event string, data interface{}) {
	r.send <- &message{
		event: event,
		data:  data,
	}
}
