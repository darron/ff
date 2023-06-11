package cmd

import "testing"

func TestCleanUpLink(t *testing.T) {
	have := "https://toronto.ctvnews.ca/2-tenants-landlord-dead-after-dispute-in-residence-near-hamilton-1.6416439?fbclid=IwAR1pZYOMAPPRwPLYtnnTbsNykwb7zqtow19jHUbA4lCxvikVy9fbnGzgmgM"
	want := "https://toronto.ctvnews.ca/2-tenants-landlord-dead-after-dispute-in-residence-near-hamilton-1.6416439"
	got := cleanUpLink(have)
	if got != want {
		t.Errorf("Want: %s Got: %s", want, got)
	}
}
