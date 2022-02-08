package currate

import (
	"net/url"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/thedustin/go-email-curator/cmd/currator/cmd"
	"github.com/thedustin/go-email-curator/config"
	"github.com/thedustin/go-email-curator/filter"
	parser "github.com/thedustin/go-gmail-query-parser"
	"github.com/thedustin/go-gmail-query-parser/criteria"
)

type MailServerCmd struct {
	MailServer   *url.URL `required:"" arg:"" env:"CURRATOR_MAIL_SERVER"`
	MailUsername string   `help:"" env:"CURRATOR_MAIL_USERNAME"`
	MailPassword string   `help:"" env:"CURRATOR_MAIL_PASSWORD"`
}

type CurrateCmd struct {
	MailServerCmd
}

var section imap.BodySectionName = func() imap.BodySectionName {
	var section imap.BodySectionName

	section.Specifier = imap.TextSpecifier
	section.Peek = true

	return section
}()

func (cmd *CurrateCmd) Run(ctx *cmd.Context) error {
	config.Get().Server = cmd.MailServer
	config.Get().Username = cmd.MailUsername
	config.Get().Password = cmd.MailPassword

	err := config.Get().FromYamlFile("./config.yaml")

	if err != nil {
		return err
	}

	ctx.Logger.Println("Create filter...")

	parser := parser.NewParser(
		criteria.ValueTransformer(messageTransformer),
		parser.FlagDefault,
	)

	filters := make([]*filter.Filter, len(config.Get().Filters))

	for i, c := range config.Get().Filters {
		f, err := filter.NewFromConfig(parser, c)

		if err != nil {
			return err
		}

		filters[i] = f
	}

	ctx.Logger.Println("Filter created")

	ctx.Logger.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(config.Get().Server.String(), nil)
	if err != nil {
		ctx.Logger.Fatal(err)
	}
	ctx.Logger.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(config.Get().Username, config.Get().Password); err != nil {
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

	msgCnt := 250

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > uint32(msgCnt) {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - uint32(msgCnt-1)
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchUid}

	messages := make(chan *imap.Message, msgCnt)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	rMessages := make([]*imap.Message, msgCnt)

	i := 1
	for msg := range messages {
		rMessages[msgCnt-i] = msg
		i++
	}

	ctx.Logger.Println("Last", msgCnt, "messages:")
	for _, msg := range rMessages {
		ctx.Logger.Printf("* [%d] %s\n", msg.Uid, msg.Envelope.Subject)

		for _, f := range filters {
			f.Perform(msg, c)
		}
	}

	if err := <-done; err != nil {
		ctx.Logger.Fatal(err)
	}

	return nil
}
