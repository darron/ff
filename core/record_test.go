package core

import (
	"testing"

	"github.com/go-faker/faker/v4"
)

func TestRecordFaker(t *testing.T) {
	r := Record{}
	err := faker.FakeData(&r)
	if err != nil {
		t.Error(err)
	}
}
