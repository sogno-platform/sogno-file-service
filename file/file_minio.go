// SPDX-License-Identifier: Apache-2.0

package file

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type NoSuchKeyError struct {
	Message string
}

func (e *NoSuchKeyError) Error() string {
	return e.Message
}

type MinIOClient struct {
	Client *minio.Client
}

func NewMinIOClient(endpoint string) (*MinIOClient, error) {

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
	return &MinIOClient{Client: client}, err
}

func (c *MinIOClient) PutObject(bucket string, key string, content io.Reader, contentSize int64, contentType string) error {

	_, err := c.Client.PutObject(
		context.Background(),
		bucket,
		key,
		content,
		contentSize,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func (c *MinIOClient) StatObject(bucket string, key string) (minio.ObjectInfo, error) {

	info, err := c.Client.StatObject(context.Background(), bucket, key, minio.StatObjectOptions{})

	if minio.ToErrorResponse(err).Code == "NoSuchKey" {
		return info, &NoSuchKeyError{Message: err.Error()}
	} else {
		return info, err
	}
}

func (c *MinIOClient) GetObjectUrl(bucket string, key string) (*url.URL, error) {

	// Generates a presigned url which expires in 7 days (max).
	return c.Client.PresignedGetObject(context.Background(), bucket, key, time.Second*604800, make(url.Values))
}

func (c *MinIOClient) ListObjects(bucket string) (<-chan minio.ObjectInfo, error) {

	return c.Client.ListObjects(context.Background(), bucket, minio.ListObjectsOptions{
		Recursive: true,
	}), nil
}

func (c *MinIOClient) DeleteObject(bucket string, key string) error {

	return c.Client.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
}
