package sendgrid

import (
	"os"
	"testing"

	"github.com/subosito/gotenv"
)

func TestSendgridSendMail(t *testing.T) {
	gotenv.Load("../.env")
	var sg = New(&Config{
		APIKey:       os.Getenv("SG_KEY"),
		SenderMail:   "thaitanloi365@gmail.com",
		SenderName:   "Thai Tan Loi",
		BccAddresses: "Thai Tan Loi,thaitanloi365@gmail.com|Tri Luong,triluongdl@gmail.com|Ajb,Ajb@legendtruckinc.com|APS,APS@legendtruckinc.com|Kevin,kevin@calibrated.io",
	})

	sg.SendMail(SendMailParams{
		Email:       "ajb@legendtruckinc.com",
		Name:        "Loi",
		Subject:     "Test eail",
		HTMLContent: "<div>Hello world</div>",
	})
}
