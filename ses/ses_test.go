package ses

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/subosito/gotenv"
)

func TestSendMail(t *testing.T) {
	gotenv.Load("../.env")
	var ses = New(&Config{
		AWSAccessKey: os.Getenv("AWS_ACCESS_KEY"),
		AWSSecretKey: os.Getenv("AWS_SECRET_KEY"),
		AWSRegion:    "us-east-1",
		Sender:       "thaitanloi365@gmail.com",
	})

	ses.SendEmail(SendEmailParams{
		ToAddresses: []*string{
			aws.String("thaitanloi365@icloud.com"),
		},
	})
}
