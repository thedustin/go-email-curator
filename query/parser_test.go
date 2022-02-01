package query_test

import (
	"reflect"
	"testing"

	"github.com/thedustin/go-email-curator/query"
)

type testcase struct {
	Name   string
	Source string
	Err    error
	Result query.TokenList
}

func TestParser(t *testing.T) {
	ts := []testcase{
		{
			Name: "Empty Query",
		},
		{
			Name:   "Simple filter",
			Source: "from:example.org",
			Result: query.TokenList{
				query.NewToken(query.TokenField, "from"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "example.org"),
			},
		},
		{
			Name:   "Complexe filter value",
			Source: "subject:(Werbung für Treppenlifte)",
			Result: query.TokenList{
				query.NewToken(query.TokenField, "subject"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "Werbung für Treppenlifte"),
			},
		},
		{
			Name:   "Negate filter",
			Source: "-older_than:7d",
			Result: query.TokenList{
				query.NewToken(query.TokenNegate, "-"),
				query.NewToken(query.TokenField, "older_than"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "7d"),
			},
		},
		{
			Name:   "OR filter",
			Source: "older_than:7d OR larger:2M",
			Result: query.TokenList{
				query.NewToken(query.TokenField, "older_than"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "7d"),
				query.NewToken(query.TokenOr, "OR"),
				query.NewToken(query.TokenField, "larger"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "2M"),
			},
		},
		{
			Name:   "AND filter",
			Source: "older_than:7d AND larger:2M",
			Result: query.TokenList{
				query.NewToken(query.TokenField, "older_than"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "7d"),
				query.NewToken(query.TokenField, "larger"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "2M"),
			},
		},
		{
			Name:   "Normal Tuesday",
			Source: "from:(@example.org) (subject:(Werbung für Treppenlifte) OR -older_than:7d) from Lorem ipsum",
			Result: query.TokenList{
				query.NewToken(query.TokenField, "from"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "@example.org"),
				query.NewToken(query.TokenGroupStart, "("),
				query.NewToken(query.TokenField, "subject"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "Werbung für Treppenlifte"),
				query.NewToken(query.TokenOr, "OR"),
				query.NewToken(query.TokenNegate, "-"),
				query.NewToken(query.TokenField, "older_than"),
				query.NewToken(query.TokenEqual, ":"),
				query.NewToken(query.TokenFieldValue, "7d"),
				query.NewToken(query.TokenGroupEnd, ")"),
				query.NewToken(query.TokenFulltext, "from"),
				query.NewToken(query.TokenFulltext, "Lorem"),
				query.NewToken(query.TokenFulltext, "ipsum"),
			},
		},
	}

	p := query.NewParser()

	for _, tc := range ts {
		t.Run(tc.Name, func(t *testing.T) {
			err := p.Parse(tc.Source)

			if tc.Err != nil && err == nil {
				t.Errorf("Error should have been %v", tc.Err)
			}

			if tc.Err == nil && err != nil {
				t.Errorf("Error should have been nil but was %v", err)
			}

			if tc.Err != nil && err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			if !reflect.DeepEqual(tc.Result, p.Result()) {
				t.Errorf("Result does not match, expected\n\t%#v\nbut got\n\t%#v", tc.Result, p.Result())
			}

			t.Logf("%s", tc.Source)
			t.Logf("%s", p.Result())
			t.Logf("%s", p.Result().Describe())
		})
	}
}
