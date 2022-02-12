// SPDX-License-Identifier: Apache-2.0

package file

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/sogno-platform/file-service/config"
)

func minioClient() (*minio.Client, error) {
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
	return client, err
}

func putObject(key string, content io.Reader, contentSize int64, contentType string) error {

	bucket := config.GlobalConfig.MinIOBucket

	client, err := minioClient()
	if err != nil {
		return err
	}
	_, err = client.PutObject(
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

	client, err := minioClient()
	if err != nil {
		return minio.ObjectInfo{}, err
	}
	bucket := config.GlobalConfig.MinIOBucket

	return client.StatObject(context.Background(), bucket, key, minio.StatObjectOptions{})

}

func getObjectUrl(key string) (*url.URL, error) {

	client, err := minioClient()
	if err != nil {
		return nil, err
	}
	bucket := config.GlobalConfig.MinIOBucket

	// Generates a presigned url which expires in 7 days (max).
	return client.PresignedGetObject(context.Background(), bucket, key, time.Second*604800, make(url.Values))
}

func listObjects() (<-chan minio.ObjectInfo, error) {

	client, err := minioClient()
	if err != nil {
		return nil, err
	}
	bucket := config.GlobalConfig.MinIOBucket

	return client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Recursive: true,
	}), nil
}

func deleteObject(key string) error {

	client, err := minioClient()
	if err != nil {
		return err
	}
	bucket := config.GlobalConfig.MinIOBucket

	return client.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
}
