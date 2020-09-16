package twilio

import (
	"fmt"
	"os"
	"testing"

	"github.com/kevinburke/rest"
	"github.com/subosito/gotenv"
)

func TestSendSMS(t *testing.T) {
	gotenv.Load("../.env")
	var twilio = New(&Config{
		AccountSID:  os.Getenv("TWILIO_ACCOUNT_SID"),
		AuthToken:   os.Getenv("TWILIO_AUTH_TOKEN"),
		SenderPhone: os.Getenv("TWILIO_SENDER_PHONE"),
	})

	var err = twilio.SendSMS("327308788", "Test")
	switch v := err.(type) {
	case *rest.Error:
		fmt.Println("v", v.Detail, v.Title)
	}
	if e, ok := err.(*Error); ok {
		fmt.Println(e)
	}
}
