package dratini

import (
	"errors"
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
