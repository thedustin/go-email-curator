package currate

import (
	"fmt"
	"time"

	"github.com/emersion/go-imap"
	"github.com/thedustin/go-gmail-query-parser/criteria"
	"github.com/thedustin/go-gmail-query-parser/lexer"
)

func messageTransformer(field string, v interface{}) []string {
	msg, ok := v.(*imap.Message)

	if !ok {
		return []string{} // or panic
	}

	switch field {
	case lexer.FieldAfter:
		return []string{msg.Envelope.Date.Format(time.RFC3339)}
	case lexer.FieldBcc:
		return addresses(msg.Envelope.Bcc)
	case lexer.FieldBefore:
		return []string{msg.Envelope.Date.Format(time.RFC3339)}
	case lexer.FieldCategory:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldCc:
		return addresses(msg.Envelope.Cc)
	case lexer.FieldDeliveredTo:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldFilename:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldFrom:
		return addresses(msg.Envelope.From)
	case lexer.FieldHas:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldIn:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldIs:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldLabel:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldLarger:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldList:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldNewerThan:
		return []string{msg.Envelope.Date.Format(time.RFC3339)}
	case lexer.FieldNewer:
		return []string{msg.Envelope.Date.Format(time.RFC3339)}
	case lexer.FieldOlderThan:
		return []string{msg.Envelope.Date.Format(time.RFC3339)}
	case lexer.FieldOlder:
		return []string{msg.Envelope.Date.Format(time.RFC3339)}
	case lexer.FieldRfc822msgid:
		return []string{msg.Envelope.MessageId}
	case lexer.FieldSize:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldSmaller:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	case lexer.FieldSubject:
		return []string{msg.Envelope.Subject}
	case lexer.FieldTo:
		return addresses(msg.Envelope.To)
	case criteria.FieldFulltext:
		panic(fmt.Sprintf("Field %q not supported yet", field))
	}

	return []string{}
}

func addresses(addr []*imap.Address) []string {
	s := make([]string, len(addr))

	for i, a := range addr {
		s[i] = address(a)
	}

	return s
}

func address(addr *imap.Address) string {
	return fmt.Sprintf("%s <%s>", addr.PersonalName, addr.Address())
}
