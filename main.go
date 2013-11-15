/*
   Copyright (C) 2013  Salsita s.r.o.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program. If not, see {http://www.gnu.org/licenses/}.
*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	auth "github.com/abbot/go-http-auth"
	collector "github.com/salsita/cider-abstract-webhook"
)

const (
	statusUnprocessableEntity = 422
	maxBodySize               = int64(10 << 20)
)

func main() {
	collector.ListenAndServe(handleGitHubHook)
}

func handleGitHubHook(w http.ResponseWriter, r *http.Request) {
	// X-Github-Event header must be present.
	event := r.Header.Get("X-Github-Event")
	if event == "" {
		http.Error(w, "X-Github-Event Header Missing", statusUnprocessableEntity)
		return
	}

	var eventBody []byte

	// Push event is different for historical reasons.
	if event == "push" {
		p := r.FormValue("payload")
		if p == "" {
			http.Error(w, "Payload Form Value Missing", statusUnprocessableEntity)
			return
		}

		eventBody = []byte(p)
	} else {
		bodyReader := http.MaxBytesReader(w, r.Body, maxBodySize)
		defer bodyReader.Close()

		body, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			http.Error(w, "Request Payload Too Large", http.StatusRequestEntityTooLarge)
			return
		}

		eventBody = body
	}

	// Unmarshal the event object.
	var eventObject map[string]interface{}
	err := json.Unmarshal(eventBody, &eventObject)
	if err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	// Publish the event.
	if err := collector.Publish("github."+event, eventObject); err != nil {
		http.Error(w, "Event Not Published", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
