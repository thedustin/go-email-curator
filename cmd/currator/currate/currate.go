package currate

import (
	"net/url"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/thedustin/go-email-curator/cmd/currator/cmd"
	"github.com/thedustin/go-email-curator/filter"
	"github.com/thedustin/go-email-curator/filter/action"
	"github.com/thedustin/go-email-curator/filter/criteria"
)

type MailServerCmd struct {
	MailServer   *url.URL `required:"" arg:"" env:"CURRATOR_MAIL_SERVER"`
	MailUsername string   `help:"" env:"CURRATOR_MAIL_USERNAME"`
	MailPassword string   `help:"" env:"CURRATOR_MAIL_PASSWORD"`
}

type CurrateCmd struct {
	MailServerCmd
}

func (cmd *CurrateCmd) Run(ctx *cmd.Context) error {
	ctx.Logger.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(cmd.MailServer.String(), nil)
	if err != nil {
		ctx.Logger.Fatal(err)
	}
	ctx.Logger.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(cmd.MailUsername, cmd.MailPassword); err != nil {
		ctx.Logger.Fatal(err)
	}
	ctx.Logger.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	ctx.Logger.Println("Mailboxes:")
	for m := range mailboxes {
		ctx.Logger.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		ctx.Logger.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		ctx.Logger.Fatal(err)
	}
	ctx.Logger.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	filters := []filter.Filter{
		filter.NewFilter(
			criteria.NewSubject("Gorillas?"),
			action.NewDelete(),
		),
	}

	ctx.Logger.Println("Last 4 messages:")
	for msg := range messages {
		ctx.Logger.Printf("* %s\n", msg.Envelope.Subject)

		for _, f := range filters {
			err := f.Execute(msg, c)

			if err != nil {
				ctx.Logger.Fatal(err)
			}
		}
	}

	if err := <-done; err != nil {
		ctx.Logger.Fatal(err)
	}

	return nil
}
