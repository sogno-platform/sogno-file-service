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

	"github.com/sogno-platform/file-service/config"
)

func minioClient() *minio.Client {
	endpoint := config.GlobalConfig.MinIOEndpoint

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

	bucket := config.GlobalConfig.MinIOBucket

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

func statObject(key string) (minio.ObjectInfo, error) {

	bucket := config.GlobalConfig.MinIOBucket

	return minioClient().StatObject(context.Background(), bucket, key, minio.StatObjectOptions{})

}

func getObjectUrl(key string) (*url.URL, error) {

	bucket := config.GlobalConfig.MinIOBucket

	// Generates a presigned url which expires in 7 days (max).
	return minioClient().PresignedGetObject(context.Background(), bucket, key, time.Second*604800, make(url.Values))
}

func listObjects() <-chan minio.ObjectInfo {

	bucket := config.GlobalConfig.MinIOBucket

	return minioClient().ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Recursive: true,
	})
}

func deleteObject(key string) error {

	bucket := config.GlobalConfig.MinIOBucket

	return minioClient().RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
}
