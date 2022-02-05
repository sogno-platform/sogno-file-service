// SPDX-License-Identifier: Apache-2.0

package file

import (
	"context"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func minioClient() *minio.Client {
	// TODO: Include configuration
	endpoint := "s3.amazonaws.com"

	// TODO: Chain credentials providers
	creds := credentials.NewEnvAWS()

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

	client := minioClient()

	bucket := "TODO: get from config"

	_, err := client.PutObject(
		context.Background(),
		bucket,
		key,
		content,
		contentSize,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}
