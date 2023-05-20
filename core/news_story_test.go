package core

import (
	"testing"
)

func TestFakeNewsStory(t *testing.T) {
	j := FakeNewsStoryJSON()
	if j == "" {
		t.Error("JSON was blank")
	}
}
