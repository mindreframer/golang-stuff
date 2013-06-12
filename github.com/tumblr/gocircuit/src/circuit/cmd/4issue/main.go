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

/*
4issue is a command-line interface to the sample issue notification system.

	% 4issue ls

Lists current unresolved issues.

	% 4issue resolve {IssueID}

Marks the issue as resolved.

	% 4issue subscribers

List all emails subscribed to receive issue notifications.

	% 4issue subscribe {Email}

Subscribe the given email to receive new issue notifications.

	% 4issue unsubscribe {Email}

Unsubscribe the given email from issue notifications.

*/
package main

import (
	_ "circuit/load"
	"circuit/use/issuefs"
	"fmt"
	"os"
)

func usage() {
	println("Usage:", os.Args[0], "(ls | resolve ID | subscribers | subscribe Email | unsubscribe Email)")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "subscribe":
		if len(os.Args) != 3 {
			usage()
		}
		if err := issuefs.Subscribe(os.Args[2]); err != nil {
			println("Email already subscribed")
			os.Exit(1)
		}
	case "subscribers":
		subs, err := issuefs.Subscribers()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Problem reading subscribers from Zookeeper (%s)\n", err)
			os.Exit(1)
		}
		for _, s := range subs {
			fmt.Printf("%s\n", s)
		}
	case "unsubscribe":
		if len(os.Args) != 3 {
			usage()
		}
		if err := issuefs.Unsubscribe(os.Args[2]); err != nil {
			println("Email not subscribed")
			os.Exit(1)
		}
	case "ls":
		issues := issuefs.List()
		for _, i := range issues {
			fmt.Printf("%s\n", i.String())
		}
	case "resolve":
		if len(os.Args) != 3 {
			usage()
		}
		id, err := issuefs.ParseID(os.Args[2])
		if err != nil {
			println("Issue ID did not parse correctly")
			os.Exit(1)
		}
		if err = issuefs.Resolve(id); err != nil {
			println("No issue with this id")
			os.Exit(1)
		}
	default:
		usage()
	}
}
