package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient MinIO 存储客户端
type MinIOClient struct {
	Client     *minio.Client
	BucketName string
}

// MinIOConfig MinIO 配置
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// NewMinIOClient 创建 MinIO 客户端
func NewMinIOClient(cfg MinIOConfig) (*MinIOClient, error) {
	// 初始化 MinIO 客户端
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure:       cfg.UseSSL,
		BucketLookup: minio.BucketLookupPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	client := &MinIOClient{
		Client:     minioClient,
		BucketName: cfg.Bucket,
	}

	// 确保 bucket 存在
	if err := client.ensureBucket(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket: %w", err)
	}

	log.Printf("MinIO client initialized: endpoint=%s, bucket=%s", cfg.Endpoint, cfg.Bucket)
	return client, nil
}

// ensureBucket 确保 bucket 存在，不存在则创建
func (c *MinIOClient) ensureBucket() error {
	ctx := context.Background()
	exists, err := c.Client.BucketExists(ctx, c.BucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = c.Client.MakeBucket(ctx, c.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("MinIO bucket created: %s", c.BucketName)
	}

	return nil
}

// UploadFile 上传文件到 MinIO
func (c *MinIOClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := c.Client.PutObject(ctx, c.BucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

// DownloadFile 从 MinIO 下载文件
func (c *MinIOClient) DownloadFile(ctx context.Context, objectName string) (io.Reader, error) {
	reader, err := c.Client.GetObject(ctx, c.BucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	return reader, nil
}

// DeleteFile 从 MinIO 删除文件
func (c *MinIOClient) DeleteFile(ctx context.Context, objectName string) error {
	err := c.Client.RemoveObject(ctx, c.BucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetPresignedURL 获取文件的预签名 URL
func (c *MinIOClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (*url.URL, error) {
	reqParams := make(url.Values)
	presignedURL, err := c.Client.PresignedGetObject(ctx, c.BucketName, objectName, expiry, reqParams)
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned url: %w", err)
	}
	return presignedURL, nil
}

// GetObjectInfo 获取文件元信息
func (c *MinIOClient) GetObjectInfo(ctx context.Context, objectName string) (minio.ObjectInfo, error) {
	objectInfo, err := c.Client.StatObject(ctx, c.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("failed to get object info: %w", err)
	}
	return objectInfo, nil
}

// CheckFileExists 检查文件是否存在
func (c *MinIOClient) CheckFileExists(ctx context.Context, objectName string) (bool, error) {
	_, err := c.Client.StatObject(ctx, c.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		// 如果 StatObject 返回错误，检查是否是 "对象不存在" 类型的错误
		// MinIO 客户端会返回包含 "NoSuchKey" 或 "NoSuchObject" 的错误信息
		errMsg := err.Error()
		if contains(errMsg, "NoSuchKey") || contains(errMsg, "NoSuchObject") || contains(errMsg, "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

// searchString 在字符串中搜索子串
func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
