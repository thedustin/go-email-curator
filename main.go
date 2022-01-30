package main

import (
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/thedustin/go-email-curator/filter"
	"github.com/thedustin/go-email-curator/filter/action"
	"github.com/thedustin/go-email-curator/filter/criteria"
)

const (
	MailServer   = ""
	MailUsername = ""
	MailPassword = ""
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(MailServer, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(MailUsername, MailPassword); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

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

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Printf("* %s\n", msg.Envelope.Subject)

		for _, f := range filters {
			err := f.Execute(msg, c)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")

	criteria.NewAnd()
}
