package query_test

import (
	"testing"

	"github.com/thedustin/go-email-curator/query"
)

type testcase struct {
	Name   string
	Source string
	Err    error
}

func TestParser(t *testing.T) {
	ts := []testcase{
		{
			Name: "Empty Query",
		},
		{
			Name:   "Normal Tuesday",
			Source: "from:(@example.org) (subject:(Werbung für Treppenlifte) OR older_than:7d)",
		},
	}

	p := query.NewParser()

	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := p.Parse(tc.Source)

			if tc.Err != nil && err == nil {
				t.Errorf("Error should have been %v", tc.Err)
			}

			if tc.Err == nil && err != nil {
				t.Errorf("Error should have been nil but was %v", err)
			}

			if tc.Err != nil && err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			t.Logf("ts: %v", actual)
		})
	}
}
