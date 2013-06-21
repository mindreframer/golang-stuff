package main

/*
  Control

  This is the part of Ground Control that handles invoking
  shell commands and exposing the UI.

  The invoking part needs commands to be present in the configuration
  like so:

  ```
  "controls" : {
    "xbmc": {
      "on" : "/etc/init.d/xbmc start",
      "off" : "/etc/init.d/xbmc stop"
    }
  }
  ```
  While controls are accessed RESTfully like so:

  ```
  POST controls/control_name/on
  POST controls/control_name/off
  POST controls/control_name/once
  GET controls/control_name/status
  ```
*/

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type Control struct {
	Mount    string
	controls map[string]interface{}
}

type ActionResult struct {
	Output string `json:"output"`
}

func NewControl(controls map[string]interface{}) (h *Control) {
	c := &Control{Mount: "/controls/", controls: controls}
	return c
}

func (self *Control) Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, self.Mount)
	hdr := w.Header()
	hdr.Add("Access-Control-Allow-Origin", "*")

	switch {
	case r.Method == "GET" && path == "all":
		enc := json.NewEncoder(w)
		enc.Encode(&self.controls)

	case r.Method == "POST":
		tuple := strings.Split(path, "/")
		if len(tuple) != 2 {
			notFound(w)
			return
		}

		control := tuple[0]
		action := tuple[1]

		cmd, err := multimap(self.controls, control, action)

		if err != nil {
			notFound(w)
			return
		}

		log.Println("Running", cmd)

		// Some commands just hang on Output(), leave it out for now
		// lets go with "ok" for every successful command.
		err = exec.Command("sh", "-c", cmd).Start()
		out := "ok"

		log.Println("Done")

		if err != nil {
			// minor error?
			w.WriteHeader(400)
			actionResponse(w, "error.")

			log.Println(err)
		}

		actionResponse(w, string(out))
	default:
		notFound(w)
	}
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	actionResponse(w, "not found")
}
func actionResponse(w http.ResponseWriter, out string) {
	enc := json.NewEncoder(w)
	enc.Encode(&ActionResult{Output: string(out)})
}

func multimap(mmap map[string]interface{}, control string, action string) (res string, err error) {
	l1 := mmap[control]
	if l1 != nil {
		l2 := l1.(map[string]interface{})[action]
		if l2 != nil {
			return l2.(string), nil
		}
	}

	return "", errors.New("No such control/action")
}
