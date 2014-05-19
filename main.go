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
	"github.com/meeko-contrib/meeko-collector-github/handler"

	"github.com/meeko-contrib/go-meeko-webhook-receiver/receiver"
	"github.com/meeko/go-meeko/agent"
)

func main() {
	receiver.ListenAndServe(&handler.GitHubWebhookHandler{
		func(eventType string, eventObject interface{}) error {
			agent.Logging.Infof("Forwarding %s", eventType)
			return agent.PubSub.Publish(eventType, eventObject)
		},
	})
}
