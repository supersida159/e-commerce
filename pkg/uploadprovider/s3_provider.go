package uploadprovider

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/supersida159/e-commerce/common"
)

type s3Provider struct {
	bucketname string
	region     string
	apikey     string
	secretkey  string
	domain     string
	session    *session.Session
}

func NewS3Provider(bucketname, region, apikey, secretkey, domain string) *s3Provider {
	provider := &s3Provider{
		bucketname: bucketname,
		region:     region,
		apikey:     apikey,
		secretkey:  secretkey,
		domain:     domain,
	}

	s3Session, err := session.NewSession(&aws.Config{
		Region: aws.String(provider.region),
		Credentials: credentials.NewStaticCredentials(
			provider.apikey,
			provider.secretkey,
			""),
	})
	if err != nil {
		log.Fatal(err)
	}
	provider.session = s3Session

	return provider

}

func (privider *s3Provider) SaveFileUploaded(ctx context.Context, data []byte, dst string) (*common.Image, error) {
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)

	_, err := s3.New(privider.session).PutObject(&s3.PutObjectInput{

		Bucket:      aws.String(privider.bucketname),
		Key:         aws.String(dst),
		ACL:         aws.String("private"),
		Body:        fileBytes,
		ContentType: aws.String(fileType),
	})

	if err != nil {
		return nil, err
	}
	img := &common.Image{
		Url:       fmt.Sprintf("%s/%s", privider.domain, dst),
		CloudName: "s3",
	}
	return img, nil
}
