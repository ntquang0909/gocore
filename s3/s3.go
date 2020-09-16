package s3

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// S3 instance
type S3 struct {
	config *Config
	logger Logger
}

// Config config
type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Logger    Logger
}

// New init instance
func New(config *Config) *S3 {
	var s3 = &S3{
		config: config,
		logger: log.New(os.Stdout, "\r\n", 0),
	}
	if config.Logger != nil {
		s3.logger = config.Logger
	}

	return s3
}

// NewSession new s3 session
func (s3 *S3) NewSession() (*session.Session, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3.config.Region),
		Credentials: credentials.NewStaticCredentials(s3.config.AccessKey, s3.config.SecretKey, ""),
	})

	return session, err
}
