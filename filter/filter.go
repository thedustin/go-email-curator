package filter

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/thedustin/go-email-curator/action"
	"github.com/thedustin/go-email-curator/config"
	parser "github.com/thedustin/go-gmail-query-parser"
	"github.com/thedustin/go-gmail-query-parser/criteria"
)

type Filter struct {
	queryStr   string
	actionList []string

	criteria criteria.Criteria
	actions  []action.Action
}

func NewFromConfig(p *parser.Parser, c config.Filter) (*Filter, error) {
	f := Filter{
		queryStr:   c.Query,
		actionList: c.Actions,
	}

	crit, err := p.Parse(c.Query)

	if err != nil {
		return nil, err
	}

	f.criteria = crit

	actions := make([]action.Action, len(c.Actions))

	for i, a := range c.Actions {
		actions[i] = action.NewActionFromString(a)
	}

	f.actions = actions

	return &f, nil
}

func (f *Filter) Perform(msg *imap.Message, c *client.Client) {
	if !f.criteria.Matches(msg) {
		//log.Default().Printf("-> %q does not match %s (%s)\n", f.queryStr, msg.Envelope.To, msg.Envelope.Date)
		return
	}

	for _, a := range f.actions {
		log.Default().Printf("%#v\n", a)
		err := a.Perform(msg, c)

		if err != nil {
			log.Default().Println(a, err)
		}
	}
}
