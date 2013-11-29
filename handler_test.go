package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// Test the special push event type.
func TestGitHubWebhookHandler_HandlePushEvent(t *testing.T) {
	var (
		postType = "push"
		postBody = getRandomPayload()

		emittedType string
		emittedBody interface{}
	)

	handler := &GitHubWebhookHandler{
		Forward: func(eventType string, eventBody interface{}) error {
			emittedType = eventType
			emittedBody = eventBody

			return nil
		},
	}

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(postBody)

	Convey("Receiving a GitHub push webhook", t, func() {
		req, err := http.NewRequest("POST", "http://example.com", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header = http.Header{
			"X-Github-Event": {postType},
		}
		req.Form = url.Values{
			"payload": {buffer.String()},
		}

		rw := httptest.NewRecorder()

		handler.ServeHTTP(rw, req)

		if rw.Code != http.StatusNoContent {
			t.Fatalf("Unexpected status code returned: expected %d, received %d %s",
				http.StatusNoContent, rw.Code, rw.Body.String())
		}

		Convey("A push event with the right payload should be emitted", func() {
			So(emittedType, ShouldEqual, "github."+postType)
			So(emittedBody, ShouldResemble, postBody)
		})
	})
}

// Test some other event type, for example pull_request.
func TestGitHubWebhookHandler_HandlePullRequestEvent(t *testing.T) {
	var (
		postType = "pull_request"
		postBody = getRandomPayload()

		emittedType string
		emittedBody interface{}
	)

	handler := &GitHubWebhookHandler{
		Forward: func(eventType string, eventBody interface{}) error {
			emittedType = eventType
			emittedBody = eventBody

			return nil
		},
	}

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(postBody)

	Convey("Receiving a GitHub pull request webhook", t, func() {
		req, err := http.NewRequest("POST", "http://example.com", &buffer)
		if err != nil {
			t.Fatal(err)
		}

		req.Header = http.Header{
			"X-Github-Event": {postType},
		}

		rw := httptest.NewRecorder()

		handler.ServeHTTP(rw, req)

		if rw.Code != http.StatusNoContent {
			t.Fatalf("Unexpected status code returned: expected %d, received %d %s",
				http.StatusNoContent, rw.Code, rw.Body.String())
		}

		Convey("A pull request event with the right payload should be emitted", func() {
			So(emittedType, ShouldEqual, "github."+postType)
			So(emittedBody, ShouldResemble, postBody)
		})
	})
}

// Helpers ---------------------------------------------------------------------

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomPayload() map[string]interface{} {
	m := make(map[string]interface{}, rand.Intn(20))
	for i := 0; i < len(m); i++ {
		k := strconv.Itoa(rand.Int())
		v := rand.Int()
		m[k] = v
	}
	return m
}
