package ses

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SES ses
type SES struct {
	*Config
}

// Config config
type Config struct {
	AWSAccessKey string
	AWSSecretKey string
	AWSRegion    string
	Sender       string
}

// New new instance
func New(config *Config) *SES {
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	return &SES{
		Config: config,
	}
}

func (s *SES) newSession() (*ses.SES, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(s.AWSRegion)},
	)
	if err != nil {
		return nil, err
	}

	log.Println(s.AWSAccessKey, s.AWSSecretKey)
	svc := ses.New(session, session.Config.WithCredentials(credentials.NewStaticCredentials(s.AWSAccessKey, s.AWSSecretKey, "")))
	return svc, nil
}

// VerifyEmail verify email
func (s *SES) VerifyEmail(email string) error {
	session, err := s.newSession()
	if err != nil {
		return err
	}

	_, err = session.VerifyEmailAddress(&ses.VerifyEmailAddressInput{EmailAddress: aws.String(email)})
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			switch awserr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, awserr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, awserr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, awserr.Error())
			default:
				log.Println(awserr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
	}

	return err
}

// SendEmail send mail
func (s *SES) SendEmail(params SendEmailParams) error {
	session, err := s.newSession()
	if err != nil {
		return err
	}

	var charset = "UTF-8"
	if params.CharSet != "" {
		charset = params.CharSet
	}

	var sender = s.Sender
	if params.Sender != "" {
		sender = params.Sender
	}

	var input = &ses.SendEmailInput{
		Source: aws.String(sender),
		Destination: &ses.Destination{
			BccAddresses: params.BccAddresses,
			CcAddresses:  params.CcAddresses,
			ToAddresses:  params.ToAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(params.HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(charset),
					Data:    aws.String(params.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(charset),
				Data:    aws.String(params.Subject),
			},
		},
	}

	// Attempt to send the email.
	_, err = session.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}

		return err
	}

	log.Printf("Email Sent to address: %v \n", params.ToAddresses)

	return nil

}
