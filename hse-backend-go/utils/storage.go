package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var minioBucket string

// InitStorage connects to MinIO and makes sure the bucket exists.
// Call this once from main.go at startup, after db.Connect().
func InitStorage() error {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	minioBucket = os.Getenv("MINIO_BUCKET")
	if minioBucket == "" {
		minioBucket = "vessel-documents"
	}

	// MINIO_USE_SSL=true for HTTPS endpoints (e.g. the central storage server),
	// false or unset for plain HTTP such as the local MinIO container.
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}
	minioClient = client

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, minioBucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, minioBucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}

// UploadFile stores an uploaded file under a namespaced path (e.g.
// "vessel-applications/{applicationId}/{category}/{uuid}.ext") and returns
// the object path to store in the database.
func UploadFile(fileHeader *multipart.FileHeader, pathPrefix string) (objectPath string, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectPath = fmt.Sprintf("%s/%s%s", pathPrefix, uuid.NewString(), ext)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = minioClient.PutObject(ctx, minioBucket, objectPath, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: fileHeader.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	return objectPath, nil
}

// GetFileDownloadURL returns a temporary (1 hour) pre-signed URL so the
// frontend can download/view a file without exposing MinIO credentials.
func GetFileDownloadURL(objectPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url, err := minioClient.PresignedGetObject(ctx, minioBucket, objectPath, time.Hour, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// DeleteFile removes a file from storage (e.g. when replacing an upload).
func DeleteFile(objectPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return minioClient.RemoveObject(ctx, minioBucket, objectPath, minio.RemoveObjectOptions{})
}
