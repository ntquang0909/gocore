package sendgrid

import (
	"os"
	"testing"

	"github.com/subosito/gotenv"
)

func TestSendgridSendMail(t *testing.T) {
	gotenv.Load("../.env")
	var sg = New(&Config{
		APIKey:     os.Getenv("SG_KEY"),
		SenderMail: "thaitanloi365@gmail.com",
		SenderName: "Thai Tan Loi",
	})

	sg.SendMail(SendMailParams{
		Email:       "thaitanloi365@gmail.com",
		Name:        "Loi",
		Subject:     "Test eail",
		HTMLContent: "<div>Hello world</div>",
	})
}
