package pkg

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
	"github.com/sirupsen/logrus"
)

type BucketBasics struct {
	S3Client *s3.Client
}

// ListBuckets lists the buckets in the current account.
func (basics BucketBasics) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	var err error
	var output *s3.ListBucketsOutput
	var buckets []types.Bucket
	bucketPaginator := s3.NewListBucketsPaginator(basics.S3Client, &s3.ListBucketsInput{})
	for bucketPaginator.HasMorePages() {
		output, err = bucketPaginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
				logrus.Errorf("You don't have permission to list buckets for this account.")
				err = apiErr
			} else {
				logrus.Errorf("Couldn't list buckets for your account. Here's why: %v\n", err)
			}
			break
		} else {
			buckets = append(buckets, output.Buckets...)
		}
	}
	return buckets, err
}

// BucketExists checks whether a bucket exists in the current account.
func (basics BucketBasics) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := basics.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				logrus.Errorf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				logrus.Errorf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		logrus.Errorf("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}

// CreateBucket creates a bucket with the specified name in the specified Region.
func (basics BucketBasics) CreateBucket(ctx context.Context, name string, region string) error {
	_, err := basics.S3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil {
		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) {
			logrus.Errorf("You already own bucket %s.\n", name)
			err = owned
		} else if errors.As(err, &exists) {
			logrus.Errorf("Bucket %s already exists.\n", name)
			err = exists
		}
	} else {
		err = s3.NewBucketExistsWaiter(basics.S3Client).Wait(
			ctx, &s3.HeadBucketInput{Bucket: aws.String(name)}, time.Minute)
		if err != nil {
			logrus.Errorf("Failed attempt to wait for bucket %s to exist.\n", name)
		}
	}
	return err
}

// UploadFile reads from a file and puts the data into an object in a bucket.
func (basics BucketBasics) UploadFile(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		logrus.Errorf("Couldn't open file %v to upload. Here's why: %v\n", fileName, err)
	} else {
		defer file.Close()
		_, err = basics.S3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
			Body:   file,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
				logrus.Errorf("Error while uploading object to %s. The object is too large.\n"+
					"To upload objects larger than 5GB, use the S3 console (160GB max)\n"+
					"or the multipart upload API (5TB max).", bucketName)
			} else {
				logrus.Errorf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
					fileName, bucketName, objectKey, err)
			}
		} else {
			err = s3.NewObjectExistsWaiter(basics.S3Client).Wait(
				ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectKey)}, time.Minute)
			if err != nil {
				logrus.Errorf("Failed attempt to wait for object %s to exist.\n", objectKey)
			}
		}
	}
	return err
}

// UploadLargeObject uses an upload manager to upload data to an object in a bucket.
// The upload manager breaks large data into parts and uploads the parts concurrently.
func (basics BucketBasics) UploadLargeObject(ctx context.Context, bucketName string, objectKey string, largeObject []byte) error {
	largeBuffer := bytes.NewReader(largeObject)
	var partMiBs int64 = 10
	uploader := manager.NewUploader(basics.S3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   largeBuffer,
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "EntityTooLarge" {
			logrus.Errorf("Error while uploading object to %s. The object is too large.\n"+
				"The maximum size for a multipart upload is 5TB.", bucketName)
		} else {
			logrus.Errorf("Couldn't upload large object to %v:%v. Here's why: %v\n",
				bucketName, objectKey, err)
		}
	} else {
		err = s3.NewObjectExistsWaiter(basics.S3Client).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			logrus.Errorf("Failed attempt to wait for object %s to exist.\n", objectKey)
		}
	}

	return err
}

// DownloadFile gets an object from a bucket and stores it in a local file.
func (basics BucketBasics) DownloadFile(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	result, err := basics.S3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			logrus.Errorf("Can't get object %s from bucket %s. No such key exists.\n", objectKey, bucketName)
			err = noKey
		} else {
			logrus.Errorf("Couldn't get object %v:%v. Here's why: %v\n", bucketName, objectKey, err)
		}
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		logrus.Errorf("Couldn't create file %v. Here's why: %v\n", fileName, err)
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		logrus.Errorf("Couldn't read object body from %v. Here's why: %v\n", objectKey, err)
	}
	_, err = file.Write(body)
	return err
}

// DownloadLargeObject uses a download manager to download an object from a bucket.
// The download manager gets the data in parts and writes them to a buffer until all of
// the data has been downloaded.
func (basics BucketBasics) DownloadLargeObject(ctx context.Context, bucketName string, objectKey string) ([]byte, error) {
	var partMiBs int64 = 10
	downloader := manager.NewDownloader(basics.S3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
	})
	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		logrus.Errorf("Couldn't download large object from %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return buffer.Bytes(), err
}

// CopyToFolder copies an object in a bucket to a subfolder in the same bucket.
func (basics BucketBasics) CopyToFolder(ctx context.Context, bucketName string, objectKey string, folderName string) error {
	objectDest := fmt.Sprintf("%v/%v", folderName, objectKey)
	_, err := basics.S3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(fmt.Sprintf("%v/%v", bucketName, objectKey)),
		Key:        aws.String(objectDest),
	})
	if err != nil {
		var notActive *types.ObjectNotInActiveTierError
		if errors.As(err, &notActive) {
			logrus.Errorf("Couldn't copy object %s from %s because the object isn't in the active tier.\n",
				objectKey, bucketName)
			err = notActive
		}
	} else {
		err = s3.NewObjectExistsWaiter(basics.S3Client).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: aws.String(objectDest)}, time.Minute)
		if err != nil {
			logrus.Errorf("Failed attempt to wait for object %s to exist.\n", objectDest)
		}
	}
	return err
}

// CopyToBucket copies an object in a bucket to another bucket.
func (basics BucketBasics) CopyToBucket(ctx context.Context, sourceBucket string, destinationBucket string, objectKey string) error {
	_, err := basics.S3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(destinationBucket),
		CopySource: aws.String(fmt.Sprintf("%v/%v", sourceBucket, objectKey)),
		Key:        aws.String(objectKey),
	})
	if err != nil {
		var notActive *types.ObjectNotInActiveTierError
		if errors.As(err, &notActive) {
			logrus.Errorf("Couldn't copy object %s from %s because the object isn't in the active tier.\n",
				objectKey, sourceBucket)
			err = notActive
		}
	} else {
		err = s3.NewObjectExistsWaiter(basics.S3Client).Wait(
			ctx, &s3.HeadObjectInput{Bucket: aws.String(destinationBucket), Key: aws.String(objectKey)}, time.Minute)
		if err != nil {
			logrus.Errorf("Failed attempt to wait for object %s to exist.\n", objectKey)
		}
	}
	return err
}

// ListObjects lists the objects in a bucket.
func (basics BucketBasics) ListObjects(ctx context.Context, bucketName string) ([]types.Object, error) {
	var err error
	var output *s3.ListObjectsV2Output
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	}
	var objects []types.Object
	objectPaginator := s3.NewListObjectsV2Paginator(basics.S3Client, input)
	for objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(ctx)
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				logrus.Errorf("Bucket %s does not exist.\n", bucketName)
				err = noBucket
			}
			break
		} else {
			objects = append(objects, output.Contents...)
		}
	}
	return objects, err
}

// DeleteObjects deletes a list of objects from a bucket.
func (basics BucketBasics) DeleteObjects(ctx context.Context, bucketName string, objectKeys []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range objectKeys {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	output, err := basics.S3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds, Quiet: aws.Bool(true)},
	})
	if err != nil || len(output.Errors) > 0 {
		logrus.Errorf("Error deleting objects from bucket %s.\n", bucketName)
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				logrus.Errorf("Bucket %s does not exist.\n", bucketName)
				err = noBucket
			}
		} else if len(output.Errors) > 0 {
			for _, outErr := range output.Errors {
				logrus.Errorf("%s: %s\n", *outErr.Key, *outErr.Message)
			}
			err = fmt.Errorf("%s", *output.Errors[0].Message)
		}
	} else {
		for _, delObjs := range output.Deleted {
			err = s3.NewObjectNotExistsWaiter(basics.S3Client).Wait(
				ctx, &s3.HeadObjectInput{Bucket: aws.String(bucketName), Key: delObjs.Key}, time.Minute)
			if err != nil {
				logrus.Errorf("Failed attempt to wait for object %s to be deleted.\n", *delObjs.Key)
			} else {
				logrus.Errorf("Deleted %s.\n", *delObjs.Key)
			}
		}
	}
	return err
}

// DeleteBucket deletes a bucket. The bucket must be empty or an error is returned.
func (basics BucketBasics) DeleteBucket(ctx context.Context, bucketName string) error {
	_, err := basics.S3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName)})
	if err != nil {
		var noBucket *types.NoSuchBucket
		if errors.As(err, &noBucket) {
			logrus.Errorf("Bucket %s does not exist.\n", bucketName)
			err = noBucket
		} else {
			logrus.Errorf("Couldn't delete bucket %v. Here's why: %v\n", bucketName, err)
		}
	} else {
		err = s3.NewBucketNotExistsWaiter(basics.S3Client).Wait(
			ctx, &s3.HeadBucketInput{Bucket: aws.String(bucketName)}, time.Minute)
		if err != nil {
			logrus.Errorf("Failed attempt to wait for bucket %s to be deleted.\n", bucketName)
		} else {
			logrus.Errorf("Deleted %s.\n", bucketName)
		}
	}
	return err
}

func NewBucketBasics(ctx context.Context, config *config.AWSConfiguration) *BucketBasics {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(config.Region))
	if err != nil {
		logrus.Fatal(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	return &BucketBasics{S3Client: s3Client}
}

type Presigner struct {
	bucketName    string
	presignClient *s3.PresignClient
}

func (presigner Presigner) GetObject(
	ctx context.Context, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(presigner.bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		logrus.Errorf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			presigner.bucketName, objectKey, err)
	}
	return request, err
}

// PutObject makes a presigned request that can be used to put an object in a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner Presigner) PutObject(
	ctx context.Context, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(presigner.bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})

	if err != nil {
		logrus.Errorf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			presigner.bucketName, objectKey, err)
	}
	return request, err
}

// DeleteObject makes a presigned request that can be used to delete an object from a bucket.
func (presigner Presigner) DeleteObject(ctx context.Context, objectKey string) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.presignClient.PresignDeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(presigner.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		logrus.Errorf("Couldn't get a presigned request to delete object %v. Here's why: %v\n", objectKey, err)
	}
	return request, err
}

func (presigner Presigner) PresignPostObject(ctx context.Context, objectKey string, lifetimeSecs int64) (*s3.PresignedPostRequest, error) {
	request, err := presigner.presignClient.PresignPostObject(ctx, &s3.PutObjectInput{
		Bucket:  aws.String(presigner.bucketName),
		Key:     aws.String(objectKey),
		Tagging: aws.String("access=public"),
	}, func(options *s3.PresignPostOptions) {
		options.Expires = time.Duration(lifetimeSecs) * time.Second
	})

	logrus.Info(request)
	if err != nil {
		logrus.Errorf("Couldn't get a presigned post request to put %v:%v. Here's why: %v\n", presigner.bucketName, objectKey, err)
	}
	return request, nil
}

func NewPresigner(ctx context.Context, config *config.AWSConfiguration) *Presigner {
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(config.Region))
	if err != nil {
		logrus.Fatal(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bucketBasics := BucketBasics{S3Client: s3Client}

	bucketExists, err := bucketBasics.BucketExists(ctx, config.BucketName)
	if err != nil {
		logrus.Fatal(err)
	}

	if !bucketExists {
		err = bucketBasics.CreateBucket(ctx, config.BucketName, config.Region)
		if err != nil {
			panic(err)
		} else {
			logrus.Info("Bucket created.")
		}
	}

	presignClient := s3.NewPresignClient(s3Client)
	return &Presigner{presignClient: presignClient, bucketName: config.BucketName}

}
