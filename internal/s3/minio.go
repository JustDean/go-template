package s3

import (
	"context"
	"errors"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	PoolSize        uint8
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type MinioManager struct {
	isSet bool
	c     chan *minio.Client
}

var mm MinioManager

func Connect(ctx context.Context, conf MinioConfig) (*MinioManager, error) {
	if mm.isSet {
		return nil, errors.New("minio manager is already set")
	}

	mm = MinioManager{isSet: true, c: make(chan *minio.Client, conf.PoolSize)}

	useSSL := false
	for i := 0; i < int(conf.PoolSize); i++ {
		newMinioClient, err := minio.New(conf.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
			Secure: useSSL,
		})
		if err != nil {
			return nil, err
		}
		mm.c <- newMinioClient
	}
	return &mm, nil
}

func GetClient() (*minio.Client, error) {
	if !mm.isSet {
		return nil, errors.New("minio manager is not set")
	}
	return <-mm.c, nil
}

func ReleseClient(c *minio.Client) error {
	if !mm.isSet {
		return errors.New("minio manager is not set")
	}
	mm.c <- c
	return nil
}
