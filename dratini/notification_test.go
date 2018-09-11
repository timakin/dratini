package dratini

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNotification(t *testing.T) {
	cases := []struct {
		Notification DratiniPushNotification
		Expected     error
	}{
		// positive cases
		{
			DratiniPushNotification{
				Tokens:   []string{"test token"},
				Platform: 1,
				Message:  "test message",
			},
			nil,
		},
		{
			DratiniPushNotification{
				Tokens:   []string{"test token"},
				Platform: 2,
				Message:  "test message",
			},
			nil,
		},
		{
			DratiniPushNotification{
				Tokens:     []string{"test token"},
				Platform:   1,
				Message:    "test message with identifier",
				Identifier: "identifier",
			},
			nil,
		},

		// negative cases
		{
			DratiniPushNotification{
				Tokens: []string{""},
			},
			errors.New("empty token"),
		},
		{
			DratiniPushNotification{
				Tokens:   []string{"test token"},
				Platform: 100, /* neither iOS nor Android */
			},
			errors.New("invalid platform"),
		},
		{
			DratiniPushNotification{
				Tokens:   []string{"test token"},
				Platform: 1,
				Message:  "",
			},
			errors.New("empty message"),
		},
	}

	for _, c := range cases {
		actual := validateNotification(&c.Notification)
		assert.Equal(t, actual, c.Expected)
	}
}

func TestSendResponse(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendResponse(w, "valid message", http.StatusOK)
		return
	}))
	defer s.Close()

	res, err := http.Get(s.URL)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(body), "{\"message\":\"valid message\"}\n")
}
