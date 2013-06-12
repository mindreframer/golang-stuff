// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"circuit/exp/shuttr/proto"
	"fmt"
	"os/exec"
	"tumblr/firehose"
)

func StreamFirehose(freq *firehose.Request) <-chan *createRequest {
	ch := make(chan *createRequest)
	go func() {
		conn := firehose.Redial(freq)
		for {
			q := filter(conn.Read())
			if q == nil {
				continue
			}
			println(fmt.Sprintf("CREATE blogID=%d postID=%d", q.TimelineID, q.PostID))
			ch <- &createRequest{
				Forwarded: false,
				Post:      q,
				ReturnResponse: func(err error) {
					if err != nil {
						println("Firehose->XCreatePost error:", err.Error())
						return
					}
				},
			}
		}
	}()
	return ch
}

func filter(e *firehose.Event) *proto.XCreatePost {
	if e.Activity != firehose.CreatePost {
		return nil
	}
	return &proto.XCreatePost{TimelineID: e.Post.BlogID, PostID: e.Post.ID}
}

func sendmail(recipient, subject, body string) error {
	cmd := exec.Command("sendmail", recipient)
	var w bytes.Buffer
	w.WriteString("Subject: ")
	w.WriteString(subject)
	w.WriteByte('\n')
	w.Write([]byte(body))
	cmd.Stdin = &w
	_, err := cmd.CombinedOutput()
	return err
}
