package action

import (
	"fmt"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type Action interface {
	// Perform executes the action on the server.
	// Return true to call ""
	Perform(*imap.Message, *client.Client) error
}

func NewActionFromString(s string) Action {
	parts := strings.SplitN(s, ":", 2)

	switch parts[0] {
	case "delete":
		return NewDelete()
	case "seen":
		return NewSeenAction()
	case "move":
		return NewMove(parts[1])
	}

	panic(fmt.Sprintf("Unknown action %q (:%q)", parts[0], parts[1]))
}
