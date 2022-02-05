// SPDX-License-Identifier: Apache-2.0

package file

import (
	"context"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func minioClient() *minio.Client {
	// TODO: Include configuration
	endpoint := "s3.amazonaws.com"

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvAWS{},
			&credentials.FileAWSCredentials{},
		},
	)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
		Secure: true,
	})
	if err != nil {
		// TODO: Return error
		log.Fatalln(err)
	}
	return client
}

func putObject(key string, content io.Reader, contentSize int64, contentType string) error {

	bucket := "TODO: get from config"

	_, err := minioClient().PutObject(
		context.Background(),
		bucket,
		key,
		content,
		contentSize,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func getObjectUrl(key string) (*url.URL, error) {

	bucket := "TODO: get from config"

	// Generates a presigned url which expires in 7 days (max).
	return minioClient().PresignedGetObject(context.Background(), bucket, key, time.Second*604800, make(url.Values))
}
