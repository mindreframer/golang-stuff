// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"log"
	"testing"
)

func TestWatchpoint_Caller(t *testing.T) {
	go shortCircuit(commandResume())
	CallerOffset(2).Break()
}

func commandResume() command {
	return command{Action: "resume"}
}

func shortCircuit(next command) {
	cmd := <-toBrowserChannel
	log.Printf("send to browser:%#v\n", cmd)
	log.Printf("received from browser:%#v\n", next)
	fromBrowserChannel <- next
}
