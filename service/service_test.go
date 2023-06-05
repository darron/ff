package service

import (
	"database/sql"
	"testing"

	"github.com/darron/ff/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestNullBool(t *testing.T) {
	type test struct {
		input null.Bool
		want  string
	}

	tests := []test{
		{input: null.Bool{NullBool: sql.NullBool{Bool: true, Valid: true}}, want: "Yes"},
		{input: null.Bool{NullBool: sql.NullBool{Bool: false, Valid: true}}, want: "No"},
		{input: null.Bool{NullBool: sql.NullBool{Bool: false, Valid: false}}, want: ""},
		{input: null.Bool{NullBool: sql.NullBool{Bool: true, Valid: false}}, want: ""},
	}

	for _, tc := range tests {
		got := nullbool(tc.input)
		if got != tc.want {
			t.Error("They must match")
		}
	}
}

func TestGet(t *testing.T) {
	conf, _ := config.New()
	_, err := Get(conf, "../views/*.html")
	assert.NoError(t, err)
}
