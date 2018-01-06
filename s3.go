package main

// S3 is a `Storer` interface that puts an ExportResult to the
// specified S3 bucket. Don't use your main AWS keys for this!! Create backup-only keys using IAM
import (
	minio "github.com/minio/minio-go"
)

type S3 struct {
	Region       string
	Bucket       string
	AccessKey    string
	ClientSecret string
	EndPoint     string
}

//
// Store puts an `ExportResult` struct to an S3 bucket within the specified directory
//
func (x *S3) Store(result *ExportResult, directory string) *Error {

	// Make sure we do not have any errors
	if result.Error != nil {
		return result.Error
	}

	// Upload to S3 Object Store.
	err := x.UploadObject(result.Path, directory+"/"+result.Filename())

	return makeErr(err, "")
}

//
// Upload to object store.
//
func (x *S3) UploadObject(filePath string, storePath string) error {

	// New returns an Amazon S3 compatible client object.
	minioClient, err := minio.New(x.EndPoint, x.AccessKey, x.ClientSecret, true)

	if err != nil {
		return err
	}

	// Upload file.
	_, err = minioClient.FPutObject(x.Bucket, storePath, filePath, minio.PutObjectOptions{})

	if err != nil {
		return err
	}

	// Return happy
	return nil
}

/* End File */
