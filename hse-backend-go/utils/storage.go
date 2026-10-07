package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client
var minioPublicClient *minio.Client
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

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client: %w", err)
	}
	minioClient = client

	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if publicEndpoint == "" {
		publicEndpoint = endpoint
	}
	publicClient, err := minio.New(publicEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
		Region: "us-east-1",
	})
	if err != nil {
		return fmt.Errorf("failed to create public-facing MinIO client: %w", err)
	}
	minioPublicClient = publicClient

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

func UploadFile(fileHeader *multipart.FileHeader, pathPrefix string) (objectPath string, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	objectPath = fmt.Sprintf("%s/%s%s", pathPrefix, uuid.NewString(), ext)

	// Browsers sometimes send an empty or generic Content-Type (especially
	// for merged/edited PDFs from third-party tools) - fall back to
	// guessing from the file extension so MinIO stores the right type and
	// browsers render it inline instead of forcing a download.
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" {
		if guessed := mime.TypeByExtension(ext); guessed != "" {
			contentType = guessed
		} else {
			contentType = "application/octet-stream"
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = minioClient.PutObject(ctx, minioBucket, objectPath, file, fileHeader.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return objectPath, nil
}

// GetFileDownloadURL returns a temporary (1 hour) pre-signed URL so the
// frontend can download/view a file without exposing MinIO credentials.
// Signed via minioPublicClient - see the Region comment in InitStorage for
// why this works despite the container not being able to reach that host.
func GetFileDownloadURL(objectPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	url, err := minioPublicClient.PresignedGetObject(ctx, minioBucket, objectPath, time.Hour, nil)
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

// UploadBytes stores raw bytes (e.g. a generated document) at an exact,
// caller-chosen object path - unlike UploadFile, there's no random uuid
// suffix, so re-generating the same document overwrites the previous one.
func UploadBytes(data []byte, objectPath string, contentType string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := minioClient.PutObject(ctx, minioBucket, objectPath, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// DownloadFileBytes reads an object's full contents back from MinIO - used
// when we need the raw bytes again after already uploading (e.g. to attach
// a just-generated document to an email) rather than re-generating it.
func DownloadFileBytes(objectPath string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	obj, err := minioClient.GetObject(ctx, minioBucket, objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	return io.ReadAll(obj)
}
