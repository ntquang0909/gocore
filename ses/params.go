package ses

// SendEmailParams params
type SendEmailParams struct {
	BccAddresses []*string
	CcAddresses  []*string
	ToAddresses  []*string

	Sender  string
	Subject string

	// CharSet -> default: "UTF-8"
	CharSet string

	// TextBody
	TextBody string

	// HtmlBody
	HTMLBody string
}
