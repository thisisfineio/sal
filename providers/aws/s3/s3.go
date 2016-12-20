package s3

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	service    *s3.S3
	Downloader *s3manager.Downloader
	Uploader   *s3manager.Uploader
)

type Bucket struct {
	Name string
}

func NewBucket(name string) *Bucket {
	return &Bucket{Name: name}
}

func (b *Bucket) List() (*s3.ListObjectsOutput, error) {
	input := &s3.ListObjectsInput{Bucket: &b.Name}
	return service.ListObjects(input)
}

func (b *Bucket) ListPath(path string) (*s3.ListObjectsOutput, error) {
	input := &s3.ListObjectsInput{Bucket: &b.Name, Prefix: &path}
	return service.ListObjects(input)
}

func (b *Bucket) GetObject(path string) (*Object, error) {
	input := &s3.GetObjectInput{Bucket: &b.Name, Key: &path}
	obj, err := service.GetObject(input)
	return &Object{Output: obj}, err
}

func (b *Bucket) GetObjectAndWriteToWriter(w io.WriterAt, path string) (int64, error) {
	input := &s3.GetObjectInput{Bucket: &b.Name, Key: &path}
	return Downloader.Download(w, input)
}

func (b *Bucket) PutObjectFromReader(r io.Reader, name string) (*s3manager.UploadOutput, error) {
	input := &s3manager.UploadInput{Body: r, Bucket: &b.Name,
		ACL: aws.String(s3.ObjectCannedACLBucketOwnerFullControl)}
	return Uploader.Upload(input)
}

type Object struct {
	Output *s3.GetObjectOutput
	data   []byte
}

func (o *Object) Data() ([]byte, error) {
	if o.data != nil {
		return o.data, nil
	}
	defer o.Output.Body.Close()
	data, err := ioutil.ReadAll(o.Output.Body)
	if err != nil {
		return nil, err
	}
	o.data = data
	return o.data, nil
}

var (
	md5DoesNotMatch = errors.New("MD5 Checksums do not match.")
)

func (o *Object) Validate() error {
	data, err := o.Data()
	if err != nil {
		return err
	}

	sum := fmt.Sprintf("%x", md5.Sum(data))
	fmt.Println(sum)
	fmt.Println(*o.Output.ETag)
	fmt.Println(*o.Output.ETag == sum)
	return nil
}
